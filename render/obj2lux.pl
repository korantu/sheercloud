use strict;
use warnings;

use Data::Dumper;

sub parse {
    my $filename = shift;
    -f $filename or die "Please provide input obj: $!";
    
    open my $file, "<", $filename; 

    my $options = {};

    while( my $line = <$file> ){
	chomp $line;
	if ( $line =~ /(^[a-z]+) (.*+)$/) {
	    if ( not exists $options->{$1}){
		$options->{$1} = ();
	    };
	    my $data = "$2";
	    chomp $data;
	    push @{$options->{$1}}, $data
	};
    };
    return $options;
}

sub generate{
    my $data = shift;
    $data or die "Please provide parsed data";
    my $result = "Out\n";
    return $result;
};

sub test {
    my $parsed = parse("reference/Coffe-Table.obj");
    print generate($parsed);
}

test();
