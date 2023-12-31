package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(t Template, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, r, data)
	}
}

func FAQ(t Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML // WARN: make sure the source is trusted
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes, we offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "How old are you?",
			Answer:   "Old enough to know better",
		},
		{
			Question: "What are your support hours?",
			Answer: `We have support staff answering emails 24/7,
    though response times may be a bit slower on weekends.`,
		},
		{
			Question: "How do I contact support?",
			Answer: `Email us: <a href="mailto:support@lenslocked.com">
    support@lenslocked.com</a>.`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		t.Execute(w, r, questions)
	}
}
