# **Comfort Cocoon** - Hotel Room Booking System

Comfort Cocoon is a hotel room booking system designed to provide an intuitive interface for both visitors and admins to efficiently manage room bookings.

---

## **Key Features**

### For Visitors:
- **Create Account/Login**: Create an account or log in to access the system.
- **Search Option**: Find available rooms by date and preferences.
- **Book Rooms**: Book available rooms easily from the displayed options.
- **Calendar Integration**: View booking dates with integrated calendar functionality.
- **Get Confirmation**: Receive instant confirmation for your room reservation.

### For Admins:
- **Manage Bookings**: View and manage all bookings.
- **Accept or Reject Bookings**: Approve or decline room booking requests.

---

## **Tech Stack**

### Frontend:
- HTML, CSS, Bootstrap, JavaScript

### Backend:
- Go (Golang)

### Database:
- MySQL

### Packages:
- Chi Router
- SCS Session Manager
- Justin NoSurf Package for CSRF Protection

---

## **Limitations**
- No transaction/payment integration is included yet.

---

## **Demonstration**

This project showcases a fully functional hotel management system with distinct functionalities for users and admins.

### Screenshots

#### 1. Homepage
![Homepage](./static/images/screenshoots/homepage.png)
The homepage provides an overview of available rooms and other hotel services.

#### 2. Registration Page
![Registration Page](./static/images/screenshoots/signup.png)
Visitors can create accounts by filling out the registration form.

#### 3. Login Page
![Login Page](./static/images/screenshoots/login.png)
Users can log in to manage their bookings and account details.

#### 4. Rooms Available
![Rooms Available](./static/images/screenshoots/room_general.png)
Browse available rooms to make reservations.

#### 5. Search Availability
![Search Availability](./static/images/screenshoots/search_avialablity.png)
Search for room availability on specific dates.

#### 6. Room Reservation
![Room Reservation](./static/images/screenshoots/make_reservation.png)
Choose your desired room and confirm the booking.

#### 7. Reservation Summary
![Reservation Summary](./static/images/screenshoots/reservation_summary.png)
Review booking details before final confirmation.

#### 8. Confirmation Page
![Confirmation Page](./static/images/screenshoots/confirmatation.png)
Confirm that your room has been successfully reserved.

#### 9. Admin Dashboard
![Admin Dashboard](./static/images/screenshoots/admin_dashboard.png)
Admin dashboard for managing users, rooms, and bookings.

#### 10. Admin Manages Bookings
![Admin Manages Bookings](./static/images/screenshoots/managing_bookins.png)
Admins can manage bookings, approve or reject them, and view system resources.

---

## **Download Instructions**
Download the source code for "<b>Comfort Cocoon</b>":

- **For Linux Only**
- FIRST, download the source code.
- Extract the file and copy the "Comfort Cocoon" folder.
- Open this folder in your preferred editor (e.g., <b>VS Code</b>).
- Install XAMPP and start MySQL by running the following command:

    ```bash
    sudo /opt/lampp/lampp start
    ```

- Open PHPMyAdmin (http://localhost/phpmyadmin).
- Create a database named <b>bookings</b>.
- Import the `bookings.sql.zip` file (located in the Database folder of your download) into the newly created database.
- Open a terminal, navigate to your project directory, and run these commands:

    ```bash
    go mod tidy
    ```

    ```bash
    ./run.sh
    ```

- Open your browser and navigate to http://localhost:8080/



**LOGIN DETAILS**

For user: You can register a new user and login then.

For Admin:

- Email: admin@admin.com
- Password: password

We hope you have a smooth experience running this project. If you encounter any issues or need assistance, please feel free to contact me via the email link provided in the last section.

---

## 🌟 You've Made It This Far!

Wow, that's fantastic! Thank you so much for your interest and support.

If you’ve enjoyed working with this project, a ⭐️ would be greatly appreciated. Your support motivates me to keep improving and working on more projects.

Feel free to share any advice or feedback—I’d love to hear from you!

---

**🔗 [Give a Star ⭐️](https://github.com/Soyaib10/comfort-cocoon)**

**💬 [Share Your Thoughts](mailto:soyaibzihad10@gmail.com)**

Once again, thank you for your encouragement! 😊
