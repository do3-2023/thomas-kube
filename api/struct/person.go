package person

type Person struct {
    ID              int     `json:"id"`
    LastName        string  `json:"last_name"`
    PhoneNumber     string  `json:"phone_number"`
    Location        string  `json:"location"`
}