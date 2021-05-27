# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Dbui < Formula
  desc "Interactive terminal user interface and CLI for database connections. MySQL, PostgreSQL. More to come."
  homepage "https://github.com/kenanbek/dbui"
  version "0.1.2"
  bottle :unneeded

  if OS.mac? && Hardware::CPU.intel?
    url "https://github.com/KenanBek/dbui/releases/download/v0.1.2/dbui_Darwin_x86_64.tar.gz"
    sha256 "cba11850c2516d18271d746cc58b511634207dd79700960a6a20d34eaa198628"
  end
  if OS.linux? && Hardware::CPU.intel?
    url "https://github.com/KenanBek/dbui/releases/download/v0.1.2/dbui_Linux_x86_64.tar.gz"
    sha256 "84df2980178a7a025e75304ba041ed2a35071a4b73b3bbd15143bbc08252e869"
  end

  def install
    bin.install "dbui"
  end

  test do
    system "#{bin/dbui}"
  end
end