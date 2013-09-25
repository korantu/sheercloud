use strict;
use warnings;

print "Starting rendering $ARGV[0] from a helper script\n";

for my $complete (10, 20, 40, 60, 80, 100) {
    print "Percent complete: $complete\n";
}

`cp /home/kdl/tmp/luxout.png $ARGV[0].png`
