# Room-Visa
Joke-project used for educational purposes.  
You can submit an application to someone's room!  
The room owner can approve or reject the application using the admin panel.  
If the application is approved, you will receive a ticket-style visa with a QR code.  
Show that ticket to the room owner to assure them of the reliability of your visa. The owner will scan the code and receive the validity result of your visa.
## Installation
```bash
git clone https://github.com/Tyulenb/room-visa.git
cd room-visa
mkdir config
cd config
```
Next command create .env file, you should configure it manually, however it requires to also configure docker-compose file
```bash
echo -e "ADMIN_AUTH_KEY=admin\nVISA_KEY=visa\nDB=postgres://user:password@db:5432/database\nSTORAGE=storage/\nDOMAIN=localhost\nPORT=:1234" | tee .env
cd ..
```
Start the containers
```bash
docker compose up --build
```
## Things to improve
A lot of things here can be improved
- Extension of admin panel functionality
- UI/UX
- Logging
- Web-site translation
- Cleanliness of code

