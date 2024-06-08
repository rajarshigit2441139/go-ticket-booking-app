package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserInfo struct {
	FirstName string
	LastName  string
	Email     string
	Tickets   uint
}

var Concert = "FOSSILS JHAR"

const concertTickets uint = 50

var remainingTickets uint = 50
var bookings = make([]UserInfo, 0)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", greetUsers)
	http.HandleFunc("/book", bookTickets)
	http.ListenAndServe(":8088", nil) // Change the port number here
}

func greetUsers(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, struct {
		Concert          string
		TotalTickets     uint
		RemainingTickets uint
	}{
		Concert:          Concert,
		TotalTickets:     concertTickets,
		RemainingTickets: remainingTickets,
	})
}

func bookTickets(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		userFirstName := r.FormValue("firstName")
		userLastName := r.FormValue("lastName")
		userEmail := r.FormValue("email")
		userTickets := r.FormValue("tickets")

		tickets, _ := strconv.Atoi(userTickets)
		isValidName, isValidMail, isValidTickets := userInputVal(userFirstName, userLastName, userEmail, uint(tickets))

		if isValidName && isValidMail && isValidTickets {
			remainingTickets = remainingTickets - uint(tickets)

			var userData = UserInfo{
				FirstName: userFirstName,
				LastName:  userLastName,
				Email:     userEmail,
				Tickets:   uint(tickets),
			}

			bookings = append(bookings, userData)
			sendTicket(userData)
			fmt.Fprintf(w, "You booked %v tickets on the name of %v %v. You will receive a confirmation email at %v. Thank you!\n", userTickets, userFirstName, userLastName, userEmail)

			if remainingTickets == 0 {
				fmt.Fprintln(w, "No available seats. Come for the next one.")
			}
		} else {
			if !isValidName {
				fmt.Fprintln(w, "First name or last name is too short")
			}
			if !isValidMail {
				fmt.Fprintln(w, "Email ID is not correct")
			}
			if !isValidTickets {
				fmt.Fprintln(w, "Enter valid number of tickets")
			}
			fmt.Fprintf(w, "Your input data is invalid. We have remaining %v tickets only. Please try again.\n", remainingTickets)
		}
	}
}

func userInputVal(firstName, lastName, email string, tickets uint) (bool, bool, bool) {
	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidMail := strings.Contains(email, "@") && strings.Contains(email, ".com")
	isValidTickets := tickets <= remainingTickets && tickets > 0
	return isValidName, isValidMail, isValidTickets
}

func sendTicket(user UserInfo) {
	time.Sleep(10 * time.Second)
	ticket := fmt.Sprintf("%v tickets for %v %v\n", user.Tickets, user.FirstName, user.LastName)
	fmt.Printf("Sending %v to %v.\n", ticket, user.Email)
}
