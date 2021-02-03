exeDirectory="radarSystemRaw"
configDirectory="radarSystemRaw/config"
runFileTemp="outMultipleTemp.lexe"
runFile="outMultipleTest.lexe"
configFile="vec.config"
# kill all screens
killall screen
# remove the config file
mv $configFile $configDirectory/$configFile
# rename the executable files
mv $runFileTemp $exeDirectory/$runFile
chmod +x $exeDirectory/$runFile
# restart the program
screen -S targetSession -d -m bash
screen -r targetSession -X stuff "cd $exeDirectory"$(echo -ne '\015')
screen -r targetSession -X stuff "sudo ./$runFile"$(echo -ne '\015')

sleep 1

# test using pgrep
if ! pgrep "out" > /dev/null
then
    echo "no program started"
else
    echo "program starts again"
fi