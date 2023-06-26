# Booking App
Go application for booking/reservation system with Monolithic architecture.

- Built in Go version 1.20
- Uses [chi](https://github.com/go-chi/chi) lightweight and composable router
- Uses [Alex Edward's SCS](https://github.com/alexedwards/scs) session management
- Uses [nosurf](https://www.github.com/justinas/nosurf) Cross-Site Request Forgery attacks prevention
- Uses [Bootstrap v5.3](https://getbootstrap.com/) templates
- Uses [VanillaJS Datepicker](https://github.com/mymth/vanillajs-datepicker) date input
- Uses [Notie JS library](https://github.com/jaredreich/notie) notifications
- Uses [SweetAlert2](https://sweetalert2.github.io/) popup alert messaging
- Uses [govalidator](https://github.com/asaskevich/govalidator) data check
- Uses [Soda](https://gobuffalo.io/documentation/database/soda/) CLI migrations
- Uses [pgx](https://github.com/jackc/pgx) PostgreSQL driver
- Uses [go-simple-mail](https://github.com/xhit/go-simple-mail) SMTP emails sending
- Uses [MailHog](https://github.com/mailhog/MailHog) email testing tool
- Uses [Foundation for Emails](https://get.foundation/emails.html) email template
- Uses [RoyalUI-Free-Bootstrap-Admin-Template](https://github.com/BootstrapDash/RoyalUI-Free-Bootstrap-Admin-Template) dashboard template
- Uses [simple-datatables](https://github.com/fiduswriter/Simple-DataTables) tables template

---
## Commands

- Run project (Linux):
````
./run.sh
````

- Run go test:
````
go test -v
````

- Run all go tests in subdirectories:
````
go test -v ./...
````

- Run coverage test command: 
````
go test -coverprofile=coverage.out && go tool cover -html=coverage.out
````

- Run migration tool (parameters: up|down|reset):
````
~/go/bin/soda migrate
````