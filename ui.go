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
    <title>Amelia</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
		<link rel="stylesheet" href="http://bootswatch.com/amelia/bootstrap.min.css">
  </head>
  <body>
		<nav class="navbar navbar-default" role="navigation">
			<div class="container">
				<!-- Brand and toggle get grouped for better mobile display -->
				<div class="navbar-header">
					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navbar-collapse">
						<span class="sr-only">Toggle navigation</span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="#">Amelia</a>
				</div>

				<!-- Collect the nav links, forms, and other content for toggling -->
				<div class="collapse navbar-collapse" id="navbar-collapse">
					<form class="navbar-form navbar-right" action="/login">
						<button type="submit" class="btn btn-default">Sign Up</button>
					</form>
				</div><!-- /.navbar-collapse -->
			</div><!-- /.container-fluid -->
		</nav>
		<div class="container">
			<div class="jumbotron">
				<h1 style="font-size: 128px;">Amelia</h1>
				<p style="font-size: 32px;">Amelia keeps track of your children when you aren't.</p>
				<p>
					<form action="/login">
						<button type="submit" style="font-size: 24px;" class="btn btn-primary btn-lg">Sign Up</button>
					</form>
				</p>
			</div>
		</div>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.1.1/js/bootstrap.min.js"></script>
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
    <title>Amelia</title>
    <meta http-equiv="content-type" content="text/html; charset=utf-8"/>
		<link rel="stylesheet" href="http://bootswatch.com/amelia/bootstrap.min.css">
  </head>
  <body>
		<nav class="navbar navbar-default" role="navigation">
			<div class="container">
				<!-- Brand and toggle get grouped for better mobile display -->
				<div class="navbar-header">
					<button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#navbar-collapse">
						<span class="sr-only">Toggle navigation</span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
						<span class="icon-bar"></span>
					</button>
					<a class="navbar-brand" href="#">Amelia</a>
				</div>

				<!-- Collect the nav links, forms, and other content for toggling -->
				<div class="collapse navbar-collapse" id="navbar-collapse">
					<ul class="nav navbar-nav navbar-right">
						<li><a href="/logout">Sign Out</a></li>
					</ul>
				</div><!-- /.navbar-collapse -->
			</div><!-- /.container-fluid -->
		</nav>
		<div class="container">
    <p>Hello, {{.Name}}! </p>
		{{if .AuthorizedWithMoves}}
		<div class="well">
			<p>Account is currently authorized with Moves.</p>
		</div>
		{{else}}
		<p><a href="/authorize" class="btn btn-default">Authenticate With Moves</a></p>
		{{end}}
		<h2>Phones</h2>
    {{range .PhoneEntries}}
      <p><b>{{.Parent}}</b>: {{.Phone}}</p>
    {{end}}
		<h3>Add Phone Number</h3>
    <form action="/addphone" method="POST" role="form">
			<div class="form-group">
				<label for="parent">Parent</label>
				<input type="text" class="form-control" id="parent" name="parent" placeholder="Enter parent name.">
			</div>
			<div class="form-group">
				<label for="phone">Phone</label>
				<input type="text" class="form-control" id="phone" name="phone" placeholder="+15553332222">
			</div>
			<button type="submit" class="btn btn-default">Add Phone Number</button>
    </form>
		<h3>Remove Phone Number</h3>
    <form action="/delphone" method="POST">
			<div class="form-group">
				<select name="parent" class="form-control">
					<option value=""></option>
					{{range .PhoneEntries}}
						<option value="{{.Parent}}">{{.Parent}}</option>
					{{end}}
				</select>
			</div>

			<button type="submit" class="btn btn-default">Remove Phone Number</button>
    </form>
		</div>
		<script src="//cdnjs.cloudflare.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
		<script src="//cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.1.1/js/bootstrap.min.js"></script>
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
	user := User{
		Name: u.String(),
	}

	userKey := datastore.NewKey(c, "User", u.ID, 0, nil)

	err := datastore.Get(c, userKey, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = datastore.NewQuery("PhoneEntry").Ancestor(userKey).GetAll(c, &user.PhoneEntries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = phoneTemplate.Execute(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
