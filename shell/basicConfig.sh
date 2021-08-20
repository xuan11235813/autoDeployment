exeDirectory="radarSystem"
configDirectory="radarSystem/config"
runFile="outMultiple.lexe"
configFile="vec.config"

mv $configFile $configDirectory/$configFile

cd $exeDirectory
sudo ifconfig can0 txqueuelen 1024
sudo ifconfig can1 txqueuelen 1024
sudo ifconfig can2 txqueuelen 1024
sudo ifconfig can3 txqueuelen 1024
canplayer -I canConf.log