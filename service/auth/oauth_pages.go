package auth

const confirmationPage = `
<html>
<head><title>SPACES authentication</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous"/>
    <script defer="defer" src="https://use.fontawesome.com/releases/v5.0.6/js/all.js"></script>
</head>
<body background="/static/background.jpg">
<div class="container">
    <div class="row">
        <div class="col-lg-6 col-md-8 ml-auto mr-auto pt-5">
            <div class="card border-success">
                <div class="card-body text-success">
					<h5 class="card-title">
						<i class="fas fa-handshake"></i>
						Hooray!
					</h5>
					<p>
						Your access token was successfully received from the SPACES identity server.
					</p>
					<p>
						You can now close this browser window and start using SpaceCTL.
					</p>
                </div>
                <div class="card-footer">
                    <button class="btn-lg btn-primary btn-block" onclick="window.close()">Close this window</button>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>
`

const errorPage = `
<html>
<head><title>SPACES authentication</title>
    <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" rel="stylesheet"
          integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous"/>
    <script defer="defer" src="https://use.fontawesome.com/releases/v5.0.6/js/all.js"></script>
</head>
<body background="/static/background.jpg">
<div class="container">
    <div class="row">
        <div class="col-lg-6 col-md-8 ml-auto mr-auto pt-5">
            <div class="card border-danger">
                <div class="card-body text-danger">
					<h5 class="card-title">
						<i class="fas fa-fire"></i>
						Oops!
					</h5>
					<p>
						An error occurred while retrieving your token from the SPACES identity server.
						Here's what the identity server says:
					</p>
					<pre>{{ .err }}</pre>
                </div>
            </div>
        </div>
    </div>
</div>
</body>
</html>
`