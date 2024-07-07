-- Create a keyspace
CREATE KEYSPACE url_shortener  WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1':  1};

-- Create a table
CREATE TABLE url_shortener.url_map (
    short_url text, 
    long_url text,
    created_at timestamp, 
    PRIMARY KEY (short_url)
);

-- Reverse table to get long url from short url
CREATE TABLE url_shortener.long_to_short (
    short_url text, 
    long_url text,
    PRIMARY KEY (long_url)
);

