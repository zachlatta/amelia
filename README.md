# Amelia ![Analytics](https://ga-beacon.appspot.com/UA-34529482-6/postman/readme?pixel)

Amelia send a text message to your parents with your current address every time
you change location.

Amelia was created during the [Burbank Game+Hack](http://burbankgamehack.com/)
and won second place. It also won the Mashery API prize.

#### How To Use

1. Sign up for Amelia on http://getamelia.com.
2. Connect Amelia with [Moves](http://moves-app.com/) on your phone.
3. Add parent phone numbers to your account.
4. All of the registered parents will be texted your current address every time
   you change location.

## Getting Started

Amelia runs on Google App Engine. You must set up the App Engine Go SDK as
documented
[here](https://developers.google.com/appengine/docs/go/gettingstarted/devenvironment).

### Configuration

You must create a `secrets.go` file that contains configuration for Amelia in
the home directory. Here's the contents of an example file:

```
package amelia

const (
	clientId        = "moves_client_id"
	clientSecret    = "moves_client_secret"
	twilioSid       = "twilio_sid"
	twilioAuthToken = "twilio_auth_token"
	twilioPhone     = "twilio_phone_number (with country code) ex. +15554443333"
	tomtomKey       = "tomtom_geocoding_api_key"
)
```

### Development

Run the application in development mode with:

    $ goapp serve

### Deployment

Before deploying, you must make the following adjustments:

* Change the app's name in `app.yaml` to the name of your application on App
  Engine.
* Make sure you have the `RedirectURL` in `moves.go` set to your app's URL.
* `secrets.go` must be created with valid values.

Once you've made those changes, just run the following while in the root of
your local repository:

    $ goapp deploy

If you have two-factor authentication configured on your account, you'll want
to authenticate with OAuth:

    $ goapp deploy -oauth

## License

[tl;dr](https://tldrlegal.com/license/mit-license)

The MIT License (MIT)

Copyright (c) 2014 Andrew Downing and Zach Latta

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
