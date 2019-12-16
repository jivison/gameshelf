# Welcome to Gameshelf

Built with revel, a high-productivity web framework for the [Go language](http://www.golang.org/).


### Start the web server:

    revel run gameshelf


## Phases

Development of Gameshelf is split into 4 phases:

    - **Phase 1**: Develop gameshelf in pure Revel/Go with views
    - **Phase 2**: Switch to API only, returning JSON
    - **Phase 3**: Integrate GraphQl
    - **Phase 4**: Develop Typescript React frontend

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


