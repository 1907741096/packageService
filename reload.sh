#! /bin/bash

function Echo() {
    local str=$1
    local color=$2
    if [ -z ${color} ]; then color=green; fi

    case ${color} in
        red) echo -e "\033[31m ${str} \033[0m" ;;
        green) echo -e "\033[32m ${str} \033[0m" ;;
        white) echo -e "\033[37m ${str} \033[0m" ;;
        yellow) echo -e "\033[33m ${str} \033[0m" ;;
        *) echo ${str} ;;
    esac
}

port=$1
if [ -z ${port} ]; then
    port=1097
fi

exe=`lsof -i tcp:${port} -sTCP:LISTEN | sed -n '2,2 p' | awk '{print $1}'`
pid=`lsof -i tcp:${port} -sTCP:LISTEN | sed -n '2,2 p' | awk '{print $2}'`
if [ -z ${exe} -a -z ${pid} ]; then
    Echo "noting listen the port ${port}."
    exit
fi


for i in `ls -l | sed -n '2,$ p' | awk '{print $9}'`
do
    if [ ${i} == ${exe} ]; then
        kill -USR2 ${pid}
        Echo "${exe} listen the port ${port} again."
        exit
    fi
done

Echo "the current directory have not ${exe}." red
