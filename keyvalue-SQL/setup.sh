user='postgres'
export PGPASSWORD='root' 
unamestr=$(uname)
if [[ "$unamestr" == 'Linux' ]]; then
    sudo apt-get update
    sudo apt-get -y install postgresql
elif [[ "$unamestr" == 'Darwin' ]]; then
     brew install postgresql
fi
psql -h localhost -p 5432 -U $user postgres -c "Drop Database If exists kvstore;" -c "Create Database kvstore;"
go mod tidy 
go run kvsql