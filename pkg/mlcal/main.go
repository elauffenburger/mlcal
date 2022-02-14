package mlcal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	baseUrl       string = "https://music-literacy.com/"
	gradesPageUrl string = baseUrl + "student/gradesstudent.jsp"
)

const loginResponseStatusSuccess string = "success"

const cookieJSESSIONID string = "JSESSIONID"

type Client interface {
	Get() (*Calendar, error)
}

type client struct {
	client *http.Client

	email    string
	password string

	courseID string
}

func (c *client) Get() (*Calendar, error) {
	// Login.
	err := c.login()
	if err != nil {
		return nil, err
	}

	// Grab the grades.
	grades, err := c.getGrades()
	if err != nil {
		return nil, err
	}

	return &Calendar{grades}, nil
}

type Grade struct {
	title string
	due   time.Time
}

func (c *client) getGrades() ([]Grade, error) {
	// Build the full URL with the course ID.
	fullUrl := url.URL{RawPath: gradesPageUrl}
	fullUrl.Query().Add("courseID", c.courseID)

	// Grab the HTML for the grades listing.
	req, err := http.NewRequest(http.MethodGet, fullUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Parse the HTML into a queryable document.
	resHtml, err := goquery.NewDocumentFromReader(bufio.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	// Parse out the grades.
	grades := make([]Grade, 0)
	errs := make([]error, 0)

	resHtml.Find("#gradeslistingdiv > .assignmentrow").Each(func(_ int, s *goquery.Selection) {
		title := s.Find(".rowtitle").Text()
		dueDate, err := time.Parse("", fmt.Sprintf("%s %d", s.Find(".rowduedate").Text(), time.Now().Year()))
		if err != nil {
			errs = append(errs, err)
			return
		}

		grades = append(grades, Grade{title, dueDate})
	})

	if len(errs) != 0 {
		return nil, errors.WithStack(errors.Errorf("error parsing grades: %s", errs))
	}

	return grades, nil
}

func (c *client) login() error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Create a new request.
	reqBodyRW := bufio.ReadWriter{}
	req, err := http.NewRequest(http.MethodPost, baseUrl, &reqBodyRW)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	// Serialize and write the request reqJson.
	reqJson, err := json.Marshal(&request{c.email, c.password})
	if err != nil {
		return err
	}

	_, err = reqBodyRW.Write(reqJson)
	if err != nil {
		return err
	}

	// Send the request.
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	// Read the response.
	resBody, err := ioutil.ReadAll(bufio.NewReader(res.Body))
	if err != nil {
		return err
	}

	// Serialize the response into json.
	var resJson struct {
		Message string `json:"message"`
		Url     string `json:"url"`
		Status  string `json:"status"`
	}
	err = json.Unmarshal(resBody, &resJson)
	if err != nil {
		return err
	}

	// Make sure the request was successful.
	if resJson.Status != string(loginResponseStatusSuccess) {
		return errors.WithStack(errors.Errorf("login failure: %s", string(resBody)))
	}

	// Make sure the JSESSIONID is present.
	var sessionID *string
	for _, c := range res.Cookies() {
		if c.Name != string(cookieJSESSIONID) {
			continue
		}

		*sessionID = c.Value
	}

	if sessionID == nil {
		return errors.WithStack(errors.Errorf("failed to extract JSESSIONID from login response"))
	}

	return nil
}

type Calendar struct {
	Grades []Grade
}

func NewClient(email, password, courseID string) Client {
	return &client{&http.Client{}, email, password, courseID}
}
