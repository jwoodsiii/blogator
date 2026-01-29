# Blogator: Go-powered command-line RSS feed aggregator

### Requirements
- postgresql >= @14
- go >= 1.24

### Installation
1. Clone the repostitory and navigate to its directory
`
git clone https://github.com/jwoodsiii/blogator.git
cd blogator
`

2. Install dependencies
`
brew install postgresql@14
brew install go
`

3. Install blogator
`go build`

### Configuration
Blogator expects a config file in your home directory *~/.blogatorconfig.json*. You will need to update the file with the url of your postgres database
`
{"db_url":"postgres://<username>:@localhost:<port/blogator?sslmode=disable","current_user_name":""}
`

### Usage
| Command | Description |
| ----------- | ----------- |
| login | switch currentUser to another, already created user |
| register | create user on db |
| reset    | delete all users and reset tables |
| users    | list all users, special flag for logged in user |
| agg      | begin loop scraping posts from added feeds |
| addfeed  | add feed to scrape for posts |
| feeds    | list feeds followed by current user |
| follow   | follow existing feed |
| following | list followed feeds |
| unfollow | unfollow feed |
| browse | view posts |
