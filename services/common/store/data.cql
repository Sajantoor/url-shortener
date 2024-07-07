-- Create a keyspace
CREATE KEYSPACE url_shortener  WITH replication = {'class': 'NetworkTopologyStrategy', 'datacenter1':  1};

-- Create a table
CREATE TABLE url_shortener.url_map (
    short_url text, 
    long_url text,
    created_at timestamp, 
    PRIMARY KEY (short_url)
);