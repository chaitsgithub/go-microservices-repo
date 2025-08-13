docker build -t my-mariadb .
echo " " 
echo "Starting mariadb Container"
docker run --name microservicesdb -d -p 3306:3306 my-mariadb