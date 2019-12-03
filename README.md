# Welcome to Gameshelf

Built with revel, a high-productivity web framework for the [Go language](http://www.golang.org/).


### Start the web server:

    revel run gameshelf


## Code Layout

The directory structure of a generated Revel application:

    conf/             Configuration directory
        app.conf      Main app configuration file
        routes        Routes definition file

    app/              App sources
        init.go       Interceptor registration
        controllers/  App controllers go here
        views/        Templates directory
        models/       App models go here
            init.go   Database configuration

    db/
        rambler.hjson Configuration for rambler, a migration tool
                      Written in hjson, a human readable JSON alternative

    messages/         Message files

    public/           Public static assets
        css/          CSS files
        js/           Javascript files
        images/       Image files

    tests/            Test suites


