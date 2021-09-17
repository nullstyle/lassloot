build:
    mkdir -p bin
    go build -o bin/ cmd/slootinfo.go
    go build -o bin/ cmd/slootview.go
    go build -o bin/ cmd/sloot2csv.go
    go build -o bin/ cmd/sloot2meshlab.go

export-csv: build
    mkdir -p export
    ./bin/sloot2csv -unoffset ~/icloud/housebuild/lidar/ground_subsample.las > export/ground_subsample.csv
    ./bin/sloot2csv -unoffset ~/icloud/housebuild/lidar/ground_points.las > export/ground_all.csv

export-meshlab: build
    mkdir -p export
    ./bin/sloot2meshlab ~/icloud/housebuild/lidar/ground_points.las > export/ground_all.txt

clean:
    rm -rf bin
    rm -rf export