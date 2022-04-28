#!/usr/bin/env ruby
require 'json'
require 'yaml'
require 'cgi'

file = File.read('full-emoji-list.html')
ugly = file.split('<tr>')
puts ugly.class
puts ugly.length
match_glyph = %r{<td class='chars'>(?<glyph>[^<]*)<\/td>}
match_codepoint = %r{<td class='code'><a href[^>]*>(?<codepoint>[^<]*)<}
#<td class='code'><a href='#1f1fa_1f1f8' name='1f1fa_1f1f8'>U+1F1FA U+1F1F8</a></td>
match_name = %r{<td class='name'>(?<name>[^<]*)<}
#<td class='name'>flag: Samoa</td>
match_data = %r{src='(?<image>data[^']*)'>}
#uglier = ugly.collect{|i|i.scan(/<td class='chars'>([^<]*)<\/td>|src='(data[^']*)'>/)}
#puts uglier.class
#puts uglier.length

# puts ugly[2094]
#puts uglier[2094].to_json
#puts ugly[2094].scan(match_glyph).to_yaml
#puts ugly[2094].scan(match_codepoint).to_yaml
#puts ugly[2094].scan(match_name).to_yaml
#puts ugly[2094].scan(match_data).to_yaml
match_all = Regexp.union(match_glyph,match_codepoint,match_name,match_data)
#irb(main):045:0> $ugly[44].to_enum(:scan,$match_all).map{Regexp.last_match}[3].names
#=> ["glyph", "codepoint", "name", "image"]
#irb(main):046:0> $ugly[44].to_enum(:scan,$match_all).map{Regexp.last_match}[3][:name]

uglier = {}
ugly.each do |u|
  glyph = nil
  codepoint = nil
  name = nil
  data = []
  # I usually like to think I know ruby. but this eluded me, and still vexes me
  # I got this from
  # https://stackoverflow.com/questions/80357/how-to-match-all-occurrences-of-a-regex
  u.to_enum(:scan,match_all).map{Regexp.last_match}.each do |m|
    unless m[:glyph].nil?
      glyph = m[:glyph]
      next
    end
    unless m[:codepoint].nil?
      codepoint=CGI.unescape_html(m[:codepoint])
      next
    end
    unless m[:name].nil?
      # remove the this is new indicator symbol
      name = CGI.unescape_html(m[:name].gsub('⊛ ',''))
      next
    end
    unless m[:image].nil?
      data << m[:image]
      next
    end
  end
  if name.nil?
    puts "Skip the empty entry"
    next
  end
  uglier[name] = []
  uglier[name] << codepoint
  uglier[name] << glyph
  # only now can we clean up the images
  data.collect! do |d|
    unless d =~ %r{^data:image/png;base64,}
      puts "We have malformed image data for #{name}, correct and try again"
      exit
    end
    d.gsub('data:image/png;base64,','')
  end
  uglier[name] += data
end

File.open("uglier.json","w"){|f|f.puts(JSON.pretty_generate(uglier))}
File.open("db.go","w") do |f|
  f.puts("package main")
  f.puts("var db = map[string][]string{")
  uglier.each do |name,data|
    f.puts("\"#{name}\":[]string{#{data.collect{|d|"\"#{d}\""}.join(',')}},")
  end
  f.puts("}")
end
puts %x{go fmt db.go}
#  49     "women holding hands: medium-dark skin tone, light skin tone":        []string{"E12.0", "1F469 1F3FE 200D 1F91D 200D 1F469 1F3FB", "     fully-qualified", "👩🏾<200d>🤝<200d>👩🏻"},                                                                                            
#package main
#//package emoji_db
#
#var db = map[string][]string{
#	"zzz":                                   []string{"E0.6", "1F4A4", "fully-qualified", "💤"},
#	"zombie":                                []string{"E5.0", "1F9DF", "fully-qualified", "🧟"},
#	"zipper-mouth face":                     []string{"E1.0", "1F910", "fully-qualified", "🤐"},
#	"zebra":                                 []string{"E5.0", "1F993", "fully-qualified", "🦓"},
#	"zany face":                             []string{"E5.0", "1F92A", "fully-qualified", "🤪"},
#	"yo-yo":                                 []string{"E12.0", "1FA80", "fully-qualified", "🪀"},
