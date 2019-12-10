# 根据数据表生成对应的model
安装： go get -u github.com/Shelnutt2/db2struct/cmd/db2struct

使用：  db2struct -H 127.0.0.1 --mysql_port=3306 -t feature --user root -p 123456 --gorm --target ./models/feature.go -v -d feature --struct=FeatureModel --package=models --guregu