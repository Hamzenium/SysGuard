import SwiftUI
import Combine

class ResourceMonitorViewModel: ObservableObject {
    @Published var cpuUsage: Double = 0
    @Published var memoryUsage: Double = 0
    @Published var diskUsage: Double = 0
    
    @Published var alertsEnabled: Bool = true
    
    @Published var cpuThreshold: String = ""
    @Published var memoryThreshold: String = ""
    @Published var diskThreshold: String = ""
    
    private var timer: AnyCancellable?
    
    init() {
        startFetching()
    }
    
    func startFetching() {
        timer = Timer.publish(every: 1, on: .main, in: .common)
            .autoconnect()
            .sink { _ in self.fetchUsage() }
    }
    
    func fetchUsage() {
        guard let url = URL(string: "http://127.0.0.1:8080/resource-usage") else { return }
        URLSession.shared.dataTask(with: url) { data, _, error in
            if let error = error {
                print("Error fetching usage:", error.localizedDescription)
                return
            }
            guard let data = data else { return }
            do {
                let usage = try JSONDecoder().decode(ResourceUsage.self, from: data)
                DispatchQueue.main.async {
                    self.cpuUsage = usage.cpu_usage
                    self.memoryUsage = usage.memory_usage
                    self.diskUsage = usage.disk_usage
                }
            } catch {
                print("Decoding error:", error.localizedDescription)
            }
        }.resume()
    }
    
    func toggleAlerts() {
        guard let url = URL(string: "http://127.0.0.1:8080/toggle-alerts") else { return }
        let body = ["enable_alerts": alertsEnabled]
        sendPostRequest(url: url, body: body)
    }
    
    func changeLimits() {
        guard let url = URL(string: "http://127.0.0.1:8080/limit-changer"),
              let cpu = Double(cpuThreshold),
              let mem = Double(memoryThreshold),
              let disk = Double(diskThreshold)
        else {
            print("Invalid threshold input")
            return
        }
        let body: [String: Any] = [
            "cpu_threshold": cpu,
            "memory_threshold": mem,
            "disk_threshold": disk
        ]
        sendPostRequest(url: url, body: body)
    }
    
    func shutdownServer(completion: @escaping () -> Void) {
        guard let url = URL(string: "http://127.0.0.1:8080/shutdown") else { return }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        URLSession.shared.dataTask(with: request) { _, _, error in
            if let error = error {
                print("Error shutting down:", error.localizedDescription)
            }
            DispatchQueue.main.async {
                completion()
            }
        }.resume()
    }
    
    private func sendPostRequest(url: URL, body: [String: Any]) {
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        do {
            request.httpBody = try JSONSerialization.data(withJSONObject: body)
            request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        } catch {
            print("Error serializing JSON:", error.localizedDescription)
            return
        }
        URLSession.shared.dataTask(with: request) { _, _, error in
            if let error = error {
                print("POST request error:", error.localizedDescription)
            }
        }.resume()
    }
}

struct ResourceUsage: Codable {
    let cpu_usage: Double
    let memory_usage: Double
    let disk_usage: Double
}

struct ContentView: View {
    @StateObject private var vm = ResourceMonitorViewModel()
    @State private var isShuttingDown = false
    
    var body: some View {
        VStack(alignment: .leading, spacing: 18) {
            Group {
                UsageBar(title: "CPU Usage", value: vm.cpuUsage, color: .red)
                UsageBar(title: "Memory Usage", value: vm.memoryUsage, color: .blue)
                UsageBar(title: "Disk Usage", value: vm.diskUsage, color: .purple)
            }
            
            Divider()
                .padding(.vertical, 10)
            
            Toggle("Enable Alerts", isOn: $vm.alertsEnabled)
                .onChange(of: vm.alertsEnabled) { _ in vm.toggleAlerts() }
                .toggleStyle(SwitchToggleStyle(tint: .pink))
                .padding(.bottom, 10)
            
            Text("Change Resource Limits")
                .font(.title3).bold()
                .foregroundColor(.primary)
            
            HStack(spacing: 12) {
                LimitInputField(text: $vm.cpuThreshold, placeholder: "CPU %")
                LimitInputField(text: $vm.memoryThreshold, placeholder: "Memory %")
                LimitInputField(text: $vm.diskThreshold, placeholder: "Disk %")
                
                Button(action: vm.changeLimits) {
                    Text("Apply")
                        .bold()
                        .padding(.vertical, 8)
                        .padding(.horizontal, 16)
                        .background(Color.accentColor)
                        .foregroundColor(.white)
                        .cornerRadius(8)
                        .shadow(radius: 2)
                }
                .buttonStyle(PlainButtonStyle())
                .keyboardShortcut(.defaultAction)
            }
            
            Divider()
                .padding(.vertical, 10)
            
            Button(action: {
                isShuttingDown = true
                vm.shutdownServer {
                    isShuttingDown = false
                    NSApplication.shared.terminate(nil)
                }
            }) {
                Text(isShuttingDown ? "Shutting down..." : "Shutdown Server")
                    .bold()
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(isShuttingDown ? Color.gray : Color.red)
                    .foregroundColor(.white)
                    .cornerRadius(12)
                    .shadow(radius: 3)
            }
            .disabled(isShuttingDown)
            
            Spacer()
        }
        .padding(24)
        .frame(width: 400) // Wider app content width
        .background(
            RoundedRectangle(cornerRadius: 20)
                .fill(Color(.windowBackgroundColor))
                .shadow(color: .black.opacity(0.15), radius: 12, x: 0, y: 6)
        )
        .padding()
    }
}

struct UsageBar: View {
    let title: String
    let value: Double
    let color: Color
    
    var body: some View {
        VStack(alignment: .leading, spacing: 6) {
            Text(title)
                .font(.headline)
                .foregroundColor(color)
            ZStack(alignment: .leading) {
                RoundedRectangle(cornerRadius: 6)
                    .frame(height: 14)
                    .foregroundColor(Color.gray.opacity(0.3))
                RoundedRectangle(cornerRadius: 6)
                    .frame(width: CGFloat(value / 100) * 360, height: 14) // adapt bar width
                    .foregroundColor(color)
                    .animation(.easeInOut(duration: 0.5), value: value)
            }
            Text(String(format: "%.1f%%", value))
                .font(.caption)
                .foregroundColor(.secondary)
                .padding(.leading, 4)
        }
    }
}

struct LimitInputField: View {
    @Binding var text: String
    let placeholder: String
    
    var body: some View {
        TextField(placeholder, text: $text)
            .textFieldStyle(RoundedBorderTextFieldStyle())
            .frame(width: 80)
            .multilineTextAlignment(.center)
            .font(.system(size: 14, weight: .semibold))
    }
}

@main
struct SysGuardApp: App {
    var body: some Scene {
        WindowGroup("SysGuard") {
            ContentView()
}

        .windowResizability(.contentSize)  // window resizes to content size tightly
    }
}
