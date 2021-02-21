echo "BASH VERSION:" ${BASH_VERSION}
if [ $# == 0 ]; then
        echo "usage 0: ./run.sh compile"
        echo "usage 1: ./run.sh [num of nodes]"
        echo "usage 2: ./run.sh kill"
        exit 1
fi

if [ $1 == "compile" ]; then
        echo "Compile"
        go build -o test_exporter main.go
        exit 0
fi

if [ $1 == "kill" ]; then
        echo "Kill all test exporters"
        pkill test_exporter
        exit 0
fi

echo "Start $1 test exporters"
for ((i = 0; i < $1; i++))
do
        nohup ./test_exporter 10000 "$(($i + 9700))" > /dev/null 2>&1 &
        # nohup ./test_exporter 10000 "$(($i + 9700))" &
done