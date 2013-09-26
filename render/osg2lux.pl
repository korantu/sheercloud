use strict;
use warnings;
use Carp;

use Data::Dumper;

my $osg = $ARGV[0];
my $lux = $ARGV[1];

die "Input and output should be provided" unless $osg and $lux;

open( my $IN, "<", $osg ) or die "Cannot open input: $!";
open( my $OUT, ">", $lux ) or die "Cannot open output: $!";

# Logging
sub kdllog {
    my $msg = shift;
    print "  ---  ", $msg, "\n";
}

# Output
sub kdlout {
    my $msg = shift;
    print $OUT $msg;
}

# Get a reasonable name for an OSG section
sub canonical {
    my $name = shift;
    return $1 if $name =~ /::(\S+)/ or $name =~ /^(\S+)/;
}

# Once context is collected, it should be niceified, such as:
#  Split strings into point lists
#  Make sure there is same amount of normals as vectors
#  Add indices
sub update_context {
    my $ctxt = shift;

    return {} unless exists $ctxt->{VertexData};

    for ( "NormalData", "VertexData", "TexCoordData", "ColorData" ) {
	next unless $ctxt->{$_};
	my @data = split " ", $ctxt->{$_};
	$ctxt->{$_} = \@data
    }


    my @data = @{$ctxt->{NormalData}};
    my @colordata = @{$ctxt->{ColorData}} if exists $ctxt->{ColorData};

    my $points = $#{$ctxt->{VertexData}}/3;

    for ( 1..$points ){
	push @{$ctxt->{NormalData}}, @data;
	push @{$ctxt->{ColorData}}, @colordata if exists $ctxt->{ColorData};
    }

    # Indices
    $ctxt->{Index} = [];
    for ( 1..($points-1) ){
	push @{$ctxt->{Index}}, (0, $_, $_+1);
    }    

    return $ctxt
}

# TODO no need for the global; 
# Current entry in the context
my $cursor = "unused";

# List of contexts
my @contexts = ();

# Read an OSG section
sub entry {
    my $name = canonical(shift);
    my $context = shift;

    kdllog "+$name";
    if ( $name eq "Geode" ) {
	$context = {};
	$cursor = "unused";
    }

    # Further data will go into the respective ->{$name}, if matched.
    if ( $name =~ /(NormalData)|(VertexData)|(TexCoordData)|(ColorData)|(Image)/ ) {
	$cursor = "$name";
    }

    while( my $line = <> ) {
	chomp $line;
	$line =~ s/^\s+//;
	$line =~ s/\s+$//;
	entry($line, $context) if $line =~ /{/;
	last if $line =~ /}/;

	if ( $cursor eq "Image" and $line =~ /^FileName\s"([^"]+)/) {
	    my $fullpath = "$1";
	    $context->{ $cursor } = ( $fullpath =~ /([^\\\/]+)$/ ? $1 : $fullpath) ;
	}

	if ( $name eq "Array" ) {
	    $context->{ $cursor } = exists $context->{$cursor} ? "$line $context->{$cursor}" : "$line";
	};

    };
    kdllog "-$name";
    if ( $name eq "Geode" ) {
	update_context($context);
	push @contexts, $context;
    }

};

my $known_materials = {};
my $material_name = "a";

# Construct material section
sub emit_material {
    my $ctxt = shift;
    return unless exists $ctxt->{Image};

    my $file = $ctxt->{Image};
    return if exists $known_materials->{$file};

    my $the_name = "$material_name";
    $known_materials->{$file} = "$the_name";
    $material_name++;

    <<END
Texture "$the_name-texture" "color" "imagemap"
	"string filename" ["$file"]
	"string wrap" ["repeat"]
	"float gamma" [2.200000000000000]

MakeNamedMaterial "$the_name"
	"bool multibounce" ["false"]
	"texture Kd" ["$the_name-texture"]
	"color Ks" [0.34237525 0.64237525 0.34237525]
	"float index" [0.000000000000000]
	"float uroughness" [0.250000000000000]
	"float vroughness" [0.250000000000000]
	"string type" ["glossy"]

END
}

sub emit_object {
    my $ctxt = shift;
    
    return if not exists $ctxt->{"VertexData"};
    
    my $material = "default";
    $material = $known_materials->{ $ctxt->{Image}} if exists $ctxt->{Image};

    my $P = join " ", @{$ctxt->{VertexData}};
    my $N = join " ", @{$ctxt->{NormalData}};
    my $uv = join " ", @{$ctxt->{TexCoordData}};
    my $indices = join " ", @{$ctxt->{Index}};

    <<END
AttributeBegin
	NamedMaterial "$material"
	Shape "mesh"
	      "normal N" [$N]
	      "point P" [$P]
	      "float uv" [$uv]
	      "integer triindices" [$indices]
AttributeEnd

END
}

sub make_scene {
    my $content = shift;
    <<END

#Global Information
LookAt 1600 1600 1600 0 0 0 0 0 1
Camera "perspective" "float fov" [45]

Film "fleximage"
"integer xresolution" [256] "integer yresolution" [256]
"integer haltspp" [30]

PixelFilter "mitchell" "float xwidth" [2] "float ywidth" [2]

Sampler "lowdiscrepancy" "string pixelsampler" ["lowdiscrepancy"]

#Scene Specific Information
WorldBegin

MakeNamedMaterial "default"
	"string type" ["matte"]

AttributeBegin
	CoordSysTransform "camera"
	LightSource "distant"
		"point from" [0 0 0] "point to" [1 1 1]
		"color L" [3 3 3]
AttributeEnd

AttributeBegin
	LightSource "distant"
		"point from" [0 0 1000] "point to" [0 0 0]
		"color L" [3 3 3]
AttributeEnd


$content

WorldEnd
END
}

entry("top");

# kdlout(Dumper(@contexts));

# everything combined:



my @entries = ( (map { emit_material($_) } @contexts),
		(map { emit_object($_) } @contexts) );

kdllog "Total $#contexts entries";

kdlout( make_scene( join " ", @entries ));
