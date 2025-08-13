docker build -t my-postgres .
echo " " 
echo "Starting Postgres DB Container"
docker run --name postgresdb -d -p 5432:5432 my-postgres