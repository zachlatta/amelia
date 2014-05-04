package gamehack

import (
	"fmt"
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/user"
)

const homePage = `
<!doctype html>
<html>
  <head>
    <title>Game+Hack</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
  </head>
  <body>
    <form action="/login" method="post">
      <div><input type="submit" value="Sign In"></div>
    </form>
  </body>
</html>
`

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, homePage)
}

var phoneTemplate = template.Must(template.New("phone").Parse(`
<!doctype html>
<html>
  <head>
    <title>Game+Hack</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
  </head>
  <body>
    <p>Hello, {{.Name}}! <a href="/logout">Sign Out</a></p>
    <form action="/addphone" method="POST">
      <div>Parent: <input type="text" name="parent"/></div>
      <div>Phone: <input type="text" name="phone"/></div>
      <div><input type="submit" value="Add Phone Number"></div>
    </form>
    <form action="/delphone" method="POST">
      <select name="parent">
        <option value=""></option>
        {{range .PhoneEntries}}
          <option value="{{.Parent}}">{{.Parent}}</option>
        {{end}}
      </select>
      <div><input type="submit" value="Remove Phone Number"></div>
    </form>
    {{range .PhoneEntries}}
      <p><b>{{.Parent}}</b>: {{.Phone}}</p>
    {{end}}
  </body>
</html>
`))

func phone(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}
	userInfo := User {
		Name: u.String(),
	}
	_, err := datastore.NewQuery("PhoneEntry").Ancestor(datastore.NewKey(c, "User", u.ID, 0, nil)).GetAll(c, &userInfo.PhoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = phoneTemplate.Execute(w, &userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}