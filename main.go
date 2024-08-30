package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	ID        int64   `json:"id"`
	Email     string  `json:"email"`
	Amount    int64   `json:"amount"`
	Profile   Profile `json:"profile"`
	Password  string  `json:"password"`
	Username  string  `json:"username"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy string  `json:"createdBy"`
}

type Profile struct {
	Dob        string `json:"dob"`
	Avatar     string `json:"avatar"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	StaticData string `json:"staticData"`
}

func getFilteredUsers(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "http://83.136.232.77:8091/users", nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response status: %s", resp.Status)
		http.Error(w, "Bad Response from External Service", http.StatusBadGateway)
		return
	}

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		log.Printf("Failed to decode response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	filteredUsers := FilterUsers(users)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredUsers); err != nil {
		log.Printf("Failed to encode response body: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Successful response to request")
}

func FilterUsers(users []User) []User {
	var filteredUsers []User

	for _, user := range users {
		filteredUser := user

		filteredUser.Password = ""
		filteredUser.Profile.StaticData = ""

		if filteredUser.Amount >= 50000 {
			filteredUser.Email = ""
			filteredUser.Username = ""
			filteredUser.Profile.FirstName = ""
			filteredUser.Profile.LastName = ""
			filteredUser.Profile.Avatar = ""
		}

		filteredUsers = append(filteredUsers, filteredUser)
	}

	log.Println("Filtering was successful")
	return filteredUsers
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/", getFilteredUsers)
	log.Printf("Server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
