# To install:
#   brew tap toadle/lit
#   brew install lit
#
# To remove:
#   brew uninstall lit
#   brew untap toadle/lit

class Lit < Formula
  version 'v0.1.0-alpha'
  desc "lit - a command-line based quick-launcher"
  homepage "https://github.com/toadle/lit"

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/toadle/lit/releases/download/#{version}/lit-v0.1-amd64-darwin.tar.gz"
    sha256 "b1f006e80ebd91bdc57f55f033363898d444b6fc8a8474a4932dbc8a88378c55"
  elsif  OS.mac? && Hardware::CPU.arm?
    url "https://github.com/toadle/lit/releases/download/#{version}/lit-v0.1-arm64-darwin.tar.gz"
    sha256 "3a862fee7e508819cc9fa22f4f8fb21fe1f6cb89d5e6405772f3086fefbeca5c"
  elsif OS.linux?
    url "https://github.com/cantino/mcfly/releases/download/#{version}/lit-v0.1-amd64-darwin.tar.gz"
    sha256 "ad4bd165cc008f45daf8c47e66d9d9f1b32eb67adb4a64abe0d2868f46c017d1"
  end

  def install
    bin.install "lit"
  end
end