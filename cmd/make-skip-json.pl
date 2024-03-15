#!/usr/bin/env perl

# Make a JSON list of SKIP codes from Kanjidic

use strict;
use warnings;
use FindBin '$Bin';

use Data::Kanji::Kanjidic 'parse_kanjidic';
use JSON::Create 'write_json';
my $k = parse_kanjidic ('/home/ben/data/edrdg/kanjidic');
my %skip;
for my $kanji (keys %$k) {
    $skip{$kanji} = $k->{$kanji}{P};
}
write_json ("$Bin/skip.json", \%skip, indent => 1, sort => 1);
