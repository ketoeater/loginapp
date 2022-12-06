package main

import (
    "html/template"
    "net/http"

    "github.com/go-redis/redis"
    "github.com/gorilla/mux"
    "github.com/gorilla/sessions"
    "github.com/gorilla/securecookie"
)

var (
    // Create a global template registry
    templates = template.Must(template.ParseGlob("templates/*.html"))

    // Create a global Redis client
    client = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })

    // Create a global session store
    store = sessions.NewCookieStore(
        securecookie.GenerateRandomKey(32),
        securecookie.GenerateRandomKey(32),
    )
)

// User represents a user of the web app
type User struct {
    ID       int
    Username string
    Password string
}

func main() {
    // Create a new router
    r := mux.NewRouter()

    // Serve static files from the "public" directory
    r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

    // Serve the home page
    r.HandleFunc("/", homeHandler).Methods("GET")

    // Serve the login page
    r.HandleFunc("/login", loginHandler).Methods("GET")

    // Handle login form submissions
    r.HandleFunc("/login", loginFormHandler).Methods("POST")

    // Handle logout requests
    r.HandleFunc("/logout", logoutHandler).Methods("POST")

    // Start the server
    http.ListenAndServe(":3000", r)
}

// homeHandler serves the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the user is logged in
    session, err := store.Get(r, "session")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Get the user ID from the session
    userID, ok := session.Values["user_id"]
    if !ok {
        // If the user is not logged in, redirect to the login page
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Get the user from the database using the user ID
    user, err := getUserByID(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute the home.html template and write the result to the response
    err = templates.ExecuteTemplate(w, "home.html", user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// loginHandler serves the login page
func loginHandler(w http.ResponseWriter, r
