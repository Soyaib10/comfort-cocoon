package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/Soyaib10/comfort-cocoon/internal/config"
	"github.com/Soyaib10/comfort-cocoon/internal/driver"
	"github.com/Soyaib10/comfort-cocoon/internal/forms"
	"github.com/Soyaib10/comfort-cocoon/internal/helpers"
	"github.com/Soyaib10/comfort-cocoon/internal/models"
	"github.com/Soyaib10/comfort-cocoon/internal/render"
	"github.com/Soyaib10/comfort-cocoon/internal/repository"
	"github.com/Soyaib10/comfort-cocoon/internal/repository/dbrepo"
)

// Repo the repository used by the handlers

var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewMysqlRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// add this to fix invalid data error

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
		Room:      room, // add this to fix invalid data error
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		// add these lines to fix bad data error
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap, // fixes error after invalid data
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// send notifications - first to guest
	htmlMessage := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear %s: <br>
		This is confirm your reservation from %s to %s.
`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@here.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	// m.App.MailChan <- msg

	// send notification to property owner
	// htmlMessage = fmt.Sprintf(`
	// 	<strong>Reservation Notification</strong><br>
	// 	A reservation has been made for %s from %s to %s.
	// 	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	// msg = models.MailData{
	// 	To:      "me@here.com",
	// 	From:    "me@here.com",
	// 	Subject: "Reservation Notification",
	// 	Content: htmlMessage,
	// }
	log.Println(msg)
	// m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

//JSON type structure 
type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	// I removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact renders the search availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// used to have next 6 lines
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	log.Println(err)
	//	m.App.Session.Put(r.Context(), "error", "missing url parameter")
	//	http.Redirect(w, r, "/", http.StatusSeeOther)
	//	return
	//}

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// ShowLogin shows the login screen
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logging the user in
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request){
	_ = m.App.Session.RenewToken(r.Context())
	
	err := r.ParseForm()
	if err != nil{
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid(){
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	var email, password string

	email = r.Form.Get("email")
	password = r.Form.Get("password")

	id, _, err := m.DB.Authenticate(email, password)

	if err != nil{
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

		return;
	}

	m.App.Session.Put(r.Context(), "user_id", id)

	user_info, err := m.DB.GetUserByID(id)

	user := models.User{
		ID: id,
		FirstName: user_info.FirstName,
		LastName: user_info.LastName,
		Phone: user_info.Phone,
		Email: user_info.Email,
		Password: user_info.Password,
		AccessLevel: user_info.AccessLevel,
	}
	data := make(map[string]interface{})
	data["user"] = user

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Logged in Successfully")
	if user.AccessLevel == 1{
		m.App.Session.Put(r.Context(), "is_admin", user.AccessLevel)
	}

	m.App.Session.Put(r.Context(), "user_information", user)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
		IsAdmin: user.AccessLevel == 1,
	})
}

// ShowSignup shows the Signup screen
func (m *Repository) ShowSignup(w http.ResponseWriter, r *http.Request){
	// log.Println("Pre signup works")
	render.Template(w, r, "register.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowSignup handles Signup the user
func (m *Repository) PostShowSignup(w http.ResponseWriter, r *http.Request){
	_ = m.App.Session.RenewToken(r.Context())

	// log.Println("Post Signup works")
	
	err := r.ParseForm()
	if err != nil{
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("first_name", "email", "password")
	form.IsEmail("email")
	form.MinLength("password", 6)

	var first_name, last_name, phone, email, password string

	first_name = r.Form.Get("first_name")
	last_name = r.Form.Get("last_name")
	phone = r.Form.Get("phone")
	email = r.Form.Get("email")
	password = r.Form.Get("password")

	id, _, _ := m.DB.Authenticate(email, password)

	user := models.User{
		ID: id,
		FirstName: first_name,
		LastName: last_name,
		Phone: phone,
		Email: email,
		Password: password,
		AccessLevel: 0,
	}
	data := make(map[string]interface{})
	data["user"] = user

	if !form.Valid() {
		render.Template(w, r, "register.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
		})
		return
	}

	err = m.DB.UserRegistration(first_name, last_name, phone, email, password)

	if err != nil{
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid registration credentials")
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)

		return;
	}

	// auto logged in after signup
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "user_information", user)

	m.App.Session.Put(r.Context(), "flash", "Registration Successful")
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Data: data,
		IsAdmin: false,
	})
	// http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ForgotPassword handles forgot password page in
func (m *Repository) ForgotPassword(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "forgot-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) ResetPassword(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostResetPassword(w http.ResponseWriter, r *http.Request){
	// log.Println("works reset password page...")

	err := r.ParseForm()

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	email := r.Form.Get("email")

	// log.Println(email)

	form := forms.New(r.PostForm)

	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "forgot-password.page.tmpl", &models.TemplateData{
			Form:      form,
		})
		return
	}

	isEmailInDatabase, err:= m.DB.IsEmailInDatabase(email)

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	if !isEmailInDatabase{
		m.App.Session.Put(r.Context(), "error", "Invalid Email Address!")

		render.Template(w, r, "forgot-password.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	user, err := m.DB.GetUserByEmail(email)
	// log.Println("email", email)

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["user"] = user

	m.App.Session.Put(r.Context(), "user", user)

	render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
		Data: data,
		Form: form,
	})
}

// UserProfile is page of user's profile
func (m *Repository) ChangePassword(w http.ResponseWriter, r *http.Request){
	// log.Println("change password comming...")

	user, _:= m.App.Session.Get(r.Context(), "user").(models.User)

	data := make(map[string]interface{})
	data["user"] = user

	
	err:= r.ParseForm()

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("password", "password_confirm")

	form.MinLength("password", 6)
	form.MinLength("password_confirm", 6)

	if !form.Valid() {
		render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
			Form:      form,
			Data: data,
		})
		return
	}

	newPassword := r.Form.Get("password")
	confirmNewPassword := r.Form.Get("password_confirm")
	email := user.Email

	// log.Println(newPassword)
	// log.Println(confirmNewPassword)
	// log.Println(user)
	// log.Println(email)

	if newPassword != confirmNewPassword{
		// log.Println("password does not match")
		m.App.Session.Put(r.Context(), "error", "Password doesn't match")
		render.Template(w, r, "reset-password.page.tmpl", &models.TemplateData{
			Form: form,
			
		})
		return
	}

	// log.Println(email);

	err = m.DB.ResetPassword(email, newPassword)

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Password changed successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// UserProfile is page of user's profile
func (m *Repository) UserProfile(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "user-profile.page.tmpl", &models.TemplateData{})
}

// LogOut works for logout
func (m *Repository) LogOut(w http.ResponseWriter, r *http.Request){
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard handles the admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminAllReservations shows the list of all reservations
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request){
	reservations, err:= m.DB.AllReservation()
	// fmt.Println(len(reservations))

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminNewReservations shows the list of all New reservations
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request){
	reservations, err:= m.DB.NewReservation()
	// fmt.Println(len(reservations))

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservation shows the reservation in the admin tool
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request){
	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])

	if err != nil{
		helpers.ServerError(w, err)
	}
	// log.Println(id)

	src := exploded[3]

	stringMap := make(map[string]string)

	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["month"] = month
	stringMap["year"] = year

	// get reservation from the database
	res, err:= m.DB.GetReservationByID(id)
	// fmt.Println(len(reservations))

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	data["reservation"] = res

	render.Template(w, r, "admin-reservations-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data: data,
		Form: forms.New(nil),
	})
}

func (m *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil{
		helpers.ServerError(w, err)
		return
	}
	
	exploded := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(exploded[4])

	if err != nil{
		helpers.ServerError(w, err)
	}
	// log.Println(id)

	src := exploded[3]

	stringMap := make(map[string]string)

	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(res)

	if err != nil{
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Changes saved")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminReservationsCalendar displays the reservation calendar
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all the restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				// it's a reservation
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				// it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
				// log.Println(y.StartDate.Format("2006-01-2"))
			}
		}
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
		// m.App.Session.Put(r.Context(), fmt.Sprintf("reservation_map_%d", x.ID), reservationMap)
	}

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminProcessReservation  marks a reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	err := m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminDeleteReservation deletes a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminPostReservationsCalendar handles post of reservation calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// log.Println("works")
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process block
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		// Get the block map from the session. Loop through entire map, if we have an entry in the map
		// that does not exist in our posted data, and if the restriction id > 0, then it is a block we need to
		// remove.
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			// ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to values > 0, and that are not in the form post
				// the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						// log.Println("would delete block", value)

						//delete the restriction by id
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							// log.Println(err)
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}
		}
	}

	// now handle new blocks
	for name, _ := range r.PostForm {
		// log.Println("Form has name", name)

		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomID, _ := strconv.Atoi(exploded[2])

			// log.Println("would insert block for room id", roomID,"for date", exploded[3])

			t, _ := time.Parse("2006-01-2", exploded[3])
			// insert a new block
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				// log.Println(err)
				helpers.ServerError(w, err)
				return
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
