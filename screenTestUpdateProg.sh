#!/bin/bash
mkdir -p temp
rm -rf temp/*
fileName="outMultiple.lexe"
zipName="radarSystem-master"
unzip $zipName.zip -d temp

progFile=temp/$zipName/$fileName

if test -f "$progFile"
then
	prog=$(top -n1 | grep -i out | awk '{print $2}')
	if test -z "$prog"
	then 
		echo "no such  program"

		if test -f "$zipName/$fileName"
		then
			sudo rm -f $zipName/$fileName
			sudo cp $progFile $zipName/$fileName
			sudo chmod +x $zipName/$fileName
			screenResult=$(command -v "screen")
			if test -z "$screenResult"
			then
				echo "no screen installed"
			else
				echo "start the program"
				killall screen
				screen -S targetSession -d -m bash
				screen -r targetSession -X stuff "cd $zipName"$(echo -ne '\015')
				screen -r targetSession -X stuff "sudo ./$fileName"$(echo -ne '\015')
			fi
		fi
	else
		echo "program id $prog"
		if test -f "$zipName/$fileName"
		then
			sudo rm -f $zipName/$fileName
			sudo cp $progFile $zipName/$fileName
			sudo chmod +x $zipName/$fileName
			screenResult=$(command -v "screen")
			if test -z "$screenResult"
			then
				echo "no screen installed"
			else
				echo "restart the program"
				sudo kill $prog
				killall screen
				screen -S targetSession -d -m bash
				screen -r targetSession -X stuff "cd $zipName"$(echo -ne '\015')
				screen -r targetSession -X stuff "sudo ./$fileName"$(echo -ne '\015')
			fi
		else
			sudo cp $progFile $zipName/$fileName
			sudo chmod +x $zipName/$fileName
		fi

	fi
fi

