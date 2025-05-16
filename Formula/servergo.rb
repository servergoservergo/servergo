class Servergo < Formula
  desc "一个简单的命令行工具，用于快速启动HTTP文件服务器"
  homepage "https://github.com/CC11001100/servergo"
  url "https://github.com/CC11001100/servergo.git", tag: "v0.1.0", revision: "HEAD"
  license "MIT"
  head "https://github.com/CC11001100/servergo.git", branch: "main"
  
  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"servergo"
  end

  test do
    # 检查版本输出以确认安装成功
    assert_match "ServerGo", shell_output("#{bin}/servergo version")
  end
end 