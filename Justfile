build:
    mkdir -p bin
    go build -o bin/ cmd/slootinfo.go
    go build -o bin/ cmd/slootview.go
    go build -o bin/ cmd/sloot2csv.go

export: build
    mkdir -p export
    ./bin/sloot2csv ~/icloud/housebuild/lidar/ground_subsample.las > export/ground_subsample.csv
    ./bin/sloot2csv ~/icloud/housebuild/lidar/ground_points.las > export/ground_all.csv

clean:
    rm -rf bin
    rm -rf export