## Bone Health Tracker
Engage with your Twitter followers, Build meaningful relationships - [https://bonehealthtracker.com](https://bonehealthtracker.com).

### What
Bone Health Tracker classifies the tweets of **our followers** according to our interests and helps us engage with their tweets to build meaningful relationships with our Twitter followers.

### Why
Be it 10 followers or 10 million followers, Someone follows us on Twitter because there's a common interest; But what is that?

Knowing the common interests with our followers can help build meaningful relationships with them.

### How
Here's where **Bone Health Tracker** does its magic! All we need to do is connect our Twitter account and set our favorite topics.

Tweets from our followers are classified using Machine Learning (So it's not just based on keywords) and We can choose to either automatically like the tweet which matches our topic to find them in the Likes section (or) Receive an email digest at the end of the day.

We can then choose to engage with those tweets from our followers to build meaningful relationships with them.

### Demo
[![Bone Health Tracker demo video](backup/demo/video_thumbnail.png)](https://www.youtube.com/watch?v=9ANVQBrb9OE)
Clicking the above image would open the demo video in YouTube.

### Technical How
1. A 100 randomized Twitter followers are stored in a Twitter list (The list is rotated every 24 hours). 
2. Tweets from the members of the list are retrieved every 15 minutes.
3. New tweets are classified using a BERT model using our set **topics of interest** as labels.
4. Tweets which matches our interests are stored to be sent via email digest and/or liked.
5. An email digest containing the tweets from the followers, Classified according to our interests are sent.

## Testing Setup

### Requirements
1. Go version 1.18.5 (or)  higher.
2. Python version 3.6.9 (or) higher.
3. PostgreSQL version 12 (or) higher.
4. Redis server version 7.0.4 (or) higher.

### Database
Create a postgres database and create tables as per the .sql file located in `/backup/db/Create-Tables.sql`.

Note: If a new role(user) is created for testing other than default postgres user, The user needs to be provided proper permissions to the tables(Currently commented in the .sql file).

### Machine Learning Server
Install the required Python packages from `requirements.txt`.

1. Create a folder `models` in the project directory.

        $ mkdir bonehealthtracker/models
   
2. Download the hugging face transformer model to the models directory.

        $ cd bonehealthtracker/models
        $ git lfs install
        $ git clone https://huggingface.co/facebook/bart-large-mnli

3. Run the service in the bonehealthtracker folder, 

        $ cd bonehealthtracker
        $ uvicorn main:app --host 0.0.0.0 --port 8000

Note: 
1. In production, I use gunicorn with multiple workers for greater performance as found in `/backup/server_scripts/gunicorn.service`. 

2. GPU can be used in the Machine Learning service if `device=0` used while initializing model in the `main.py`.

### Twitter Developer Setup
1. Get access to developer portal at Twitter - [https://developer.twitter.com/](https://developer.twitter.com/).
2. Create a new project `bonehealthtracker-test` and new app `bonehealthtracker-test`.
3. Type of App: `Web App`.
4. Get the `Client ID`, `Client Secret`.
5. Setup OAuth 2.0,
   
   a. Callback URI / Redirect URL:  `[https project domain]/users/connect`.

   b. Website URL: `[https project domain]`

Note: For local testing, A reverse proxy service with ssl like ngrok can be used. 
After launching the application,
        
        ./ngrok http 3000


### Configuration
The configuration file is located in `/secrets/fragmenta.json`

The following fields are mandatory for the function of the application -

1. `db_url`: URL for the postgres database server.

Note: `db`, `db_pass`, `db_user`, `db_port`, can also be optionally used instead of `db_url` and if the db user is other than `postgres`, Proper permissions for the database needs to be provided as per `/backup/db/Create-Tables.sql`.

2. `port`: Port for the application, `3000` is set by default in the configuration.
3. `root_url`: **[https project domain]** mentioned earlier in the Twitter developer portal setup.
4. `client_Id`: **Client ID** mentioned earlier in the Twitter developer portal setup.
5. `client_secret`: **Client Secret** mentioned earlier in the Twitter developer portal setup.
6. `classifier_server`: **Machine Learning Server** created earlier, `[Machine Learning Server IP]:8000/classification`.
7. `redis_server`: `[Redis server IP]:6379`.
8. `twitter_redirect_uri`: `[https project domain]/users/connect`.

### Scheduling the services

`services.go` file consists of timings for scheduling the Email Digest service.

### Building the application
After ensuring Go is installed in the system,
        
        $ cd bonehealthtracker
        $ go build bonehealthtracker

### Run the application
After ensuring that Postgres, Redis and Machine Learning server is running,

    $ ./bonehealthtracker

### Usage
1. Visit the Bone Health Tracker project domain you have set in `root_url`.
2. Register a user account (With a valid email address).
3. Connect your twitter account (With couple of followers with recent tweets) and Authorize the app.
4. Set the topics for classification.
5. **Email Digest** is required (selected by default) for receiving the classified tweets by emails and/or **Auto Like** for automatically liking those tweets to later find them in the Likes section of your Twitter profile.

