**This project is about booking hotel rooms.**
**Key features**
For Visitors:
    - Create account/Login
    - Search Option
    - Book Rooms
    - CalendarIntegration
    - Get confirmation

For Admin:
    - Manage the bookings
    - Accept or reject a bookings


**Tech Stack**
Front End:
    - HTML, CSS, Bootstrap, JavaScript
BackEnd:
    - Go
Database:
    - MySQL


**Limitatios:** 
    - No transection features

# **Demonstration:** 

This project is about a hotel management system, including user and admin functionalities.

## Screenshots



### 1. Homepage
![Homepage](./static/images/screenshoots/homepage.png)

The homepage displays an overview of available rooms and other services.

### 2. Registration Page
![Registration Page](./static/images/screenshoots/signup.png)

The registration page allows new users to sign up by providing their details.

### 3. Login Page
![Login Page](./static/images/screenshoots/login.png)

The login page where registered users can log in to access their accounts.

### 4. Rooms Available
![Rooms Available](./static/images/screenshoots/room_general.png)

This page shows the available rooms that users can book.

### 5. Search Avialablity
![Rooms Available](./static/images/screenshoots/search_avialablity.png)

This page shows the available dates that users can book.

### 6. Room Reservation
![Room Reservation](./static/images/screenshoots/make_reservation.png)

The room reservation page where users can select their preferred room and confirm the booking.

### 7. Reservation Summary
![Confirmation Page](./static/images/screenshoots/reservation_summary.png)

This page sums up reservation and displays the booking details.

### 8. Confirmation Page
![Confirmation Page](./static/images/screenshoots/confirmatation.png)

This page confirms the user's reservation and displays the booking details.

### 9. Admin Dashboard
![Admin Dashboard](./static/images/screenshoots/admin_dashboard.png)

The admin dashboard where the admin can manage users, reservations, and room availability.

### 10. Admin Manages
![Admin Manages](./static/images/screenshoots/managing_bookins.png)

This page allows the admin to manage all the resources in the system including users and rooms.


**Download Instructions**
Download Source Code "<b>Comfort Cocoon</b>"

- <b>Its all about for linux</b>
- FIRST Download the source code
- Extract the file and copy "Comfort Cocoon" folder
- Open this Comfort Cocoon folder in your fovourite Editor (Like: <b>Virtual studio code</b>)
- Install XAMPP and go to your command prompt and run

```
sudo /opt/lampp/lampp start
```

- Then Open PHPMyAdmin (http://localhost/phpmyadmin)

- Then create a Database <b>bookings</b>
- Import bookings.sql.zip file(given inside the Database folder in your download folder) in this "bookings" Database
- Open the terminal and go to your project name directory and run

```
go mod tidy
```

```
./run.sh
```

- Go to your browser and go http://localhost/8080/

**LOGIN DETAILS**

Create your own staff

For Admin

- Email: admin@admin.com
- Password: password