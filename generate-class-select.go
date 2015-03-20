// Generate Class Select fetches the attendees (classes(, facilities, staff))
// from xedule.novaember.com, forms it into a list of HTML option tags, and
// writes it to the disk as 'attendee-data.txt'.
// Optionally, the location ID (lid) can be given as an argument. The default is 34.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// attendee in the format of 'xedule.novameber.com'.
type attendee struct {
	Id   int
	Name string
	Type int
}

func main() {

	// Id 34 is Boumaboulevard.
	lid := flag.Int("lid", 34, "location id")
	flag.Parse()

	atts, err := fetch(*lid)
	if err != nil {
		log.Fatal("Error fetching attendees;", err)
		return
	}

	data := []byte(fmt.Sprintf("<!-- Attendees for location ID %d, last updated %s -->\n",
		*lid, time.Now().Format("2006-01-02")))
	classes := format(atts, 1)
	//staffs := format(atts, 2)
	//facilities := format(atts, 3)

	// Classes first, then facilities and then staff
	data = append(data, classes...)
	//data = append(data, facilities...)
	//data = append(data, staffs...)

	err = save(data)
	if err != nil {
		log.Fatal("Error writing to file;", err)
	}
}

// fetch returns the attendees for the given location id.
// Data is fetched from xedule.novaember.com
func fetch(id int) ([]attendee, error) {
	res, err := http.Get(fmt.Sprintf("http://xedule.novaember.com/attendees.%d.json", id))

	if err != nil {
		return []attendee{}, err
	}

	cont, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		log.Fatal("Error while reading from response;", err)
		return []attendee{}, err

	}

	var atts []attendee
	err = json.Unmarshal(cont, &atts)

	if err != nil {
		log.Fatal("Error while unmarshalling JSON;", err)
		return []attendee{}, err
	}

	return atts, nil
}

// format converts the given []attendee to a byte array
// as HTML option tags, ready to be inserted into a select tag.
func format(atts []attendee, ty int) []byte {
	var data []byte

	var last string
	for _, a := range atts {
		if ty > 0 && a.Type != ty {
			continue
		}

		// Not sure if you like this.
		if firstName := a.Name[:6]; last != "" && firstName != last[:6] {
			str := fmt.Sprintf("<option disabled>--%s--</option>\n", firstName)
			data = append(data, []byte(str)...)
		}

		str := fmt.Sprintf("<option value=\"%d\">%s</option>\n", a.Id, a.Name)
		data = append(data, []byte(str)...)

		last = a.Name
	}

	return data
}

// save writes the bytes to the disk in the file 'attandee-data.txt'.
func save(data []byte) error {
	return ioutil.WriteFile("attendee-data.txt", data, 0644)
}
