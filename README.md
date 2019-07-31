This is a simple osquery tls server, for use wherever.  It requires
postgresql somewhere, with a table you create yourself, and a
functioning elasticsearch for dumping the logs into.

It's functionality is driven by our business need, but features can be
requested as issues on github.com/pubkraal/ostls.

This is the table code by the way:


    CREATE TABLE host (
        id SERIAL PRIMARY KEY,
        identifier text,
        uuid uuid,
        hostname text,
        token uuid,
        enrolled timestamp with time zone
    );
    CREATE UNIQUE INDEX host_pkey ON host(id int4_ops);
    CREATE INDEX ix_token ON host(token uuid_ops);


Also, I still don't really trust running code specific webservers etc so
typically I'd just let this run and listen on localhost only, while
making sure that something like nginx proxies the requests onwards to
this thing. But you do you. In any case this entire thing will only work
with a valid TLS cert and many things. Run `./ostls -help` for
instructions.
