const dbName = process.env.MONGO_INITDB_DATABASE;
const user   = process.env.MONGO_INITDB_USERNAME;
const pwd    = process.env.MONGO_INITDB_PASSWORD;

db.createUser(
    {
        user: user,
        pwd: pwd,
        roles: [
            {
                role: "readWrite",
                db: dbName
            }
        ]
    }
);

db.init_collection.insertOne({init: true});

print(`INIT: Created user '${user}' with readWrite on '${dbName}'`);