package mlcal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	baseUrl   string = "https://music-literacy.com/"
	loginUrl  string = baseUrl + "loginprocessNew.jsp"
	gradesUrl string = baseUrl + "student/gradesProcess.jsp"
)

const loginResponseStatusSuccess string = "success"

const actionGetGrades string = "getgradesforstudent"

type Client interface {
	Get() (*Calendar, error)
	GetICS() (string, error)
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

	// Grab the assignments.
	assignments, err := c.getAssignments()
	if err != nil {
		return nil, err
	}

	return &Calendar{time.Now(), assignments}, nil
}

func (c *client) GetICS() (string, error) {
	cal, err := c.Get()
	if err != nil {
		return "", err
	}

	return cal.ToICS().Serialize(), err
}

func (c *client) getAssignments() ([]Assignment, error) {
	type request struct {
		Action   string `json:"action"`
		CourseID string `json:"course_ID"`
	}

	reqJsonBody, err := json.Marshal(&request{actionGetGrades, c.courseID})
	if err != nil {
		return nil, err
	}

	// Grab the JSON for the class grades.
	req, err := http.NewRequest(http.MethodPost, gradesUrl, bufio.NewReader(bytes.NewReader(reqJsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	resJsonBody, err := ioutil.ReadAll(bufio.NewReader(res.Body))
	if err != nil {
		return nil, err
	}

	var resBody struct {
		Assignments []struct {
			Title     string `json:"title"`
			DueDate   string `json:"duedate"`
			MaxPoints string `json:"maxpoints"`
		} `json:"Assignments"`
	}

	err = json.Unmarshal(resJsonBody, &resBody)
	if err != nil {
		return nil, err
	}

	// Transform the parsed JSON.
	assignments := make([]Assignment, 0)
	timeFormat := "2006-01-02 15:04:05.0"
	for _, a := range resBody.Assignments {
		dueDate, err := time.Parse(timeFormat, a.DueDate)
		if err != nil {
			return nil, err
		}

		maxPoints, err := strconv.Atoi(a.MaxPoints)
		if err != nil {
			return nil, err
		}

		assignments = append(assignments, Assignment{a.Title, dueDate, maxPoints})
	}

	return assignments, nil
}

func (c *client) login() error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Serialize and write the request reqJson.
	reqJson, err := json.Marshal(&request{c.email, c.password})
	if err != nil {
		return err
	}

	// Create the request.
	reqBody := bufio.NewReader(bytes.NewReader(reqJson))
	req, err := http.NewRequest(http.MethodPost, loginUrl, reqBody)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

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

	return nil
}

func NewClient(email, password, courseID string) (Client, error) {
	cookieJar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return nil, err
	}

	httpClient := http.Client{
		Jar: cookieJar,
	}

	return &client{&httpClient, email, password, courseID}, nil
}
