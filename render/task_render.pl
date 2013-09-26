my $me = `pwd`;
chomp $me;
my $base = $ARGV[0];
$base or die "Scene name is not provided";
my $place = "/tmp/render".`date +%s`;

mkdir $place;

my $scene = "$place/scene.lux";
my $picture = "$place/luxout.png";

my $renderer = "/home/kdl/lux/lux-v1.2.1-i686-sse2/luxconsole";
my $converter = "/home/kdl/git/github/sheercloud/render/osg2lux.pl";

print "Starting rendering of $file in $me\nTemporary location: $place\n";

system("perl", $converter, $base, $scene); 

system($renderer, $scene);



