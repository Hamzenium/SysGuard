// swift-tools-version:5.9
import PackageDescription

let package = Package(
    name: "MySwiftApp",
    platforms: [
        .macOS(.v13)
    ],
    products: [
        .executable(name: "MySwiftApp", targets: ["MySwiftApp"])
    ],
    targets: [
        .executableTarget(
            name: "MySwiftApp",
            dependencies: [],
            path: "Sources"
        )
    ]
)
