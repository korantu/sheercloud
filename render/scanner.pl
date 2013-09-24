use strict;
use warnings;
use File::Find;
use Getopt::Long;

my $script = "worker.pl";
my $root = ".";

my $HELP =<<HELP;
-script <full path script to run for each job>;
-root <where to scan>
HELP

GetOptions ("script=s" => \$script,
	    "root=s" => \$root,   
    ) or die $HELP;


sub doRender {
    my $filename = "$_";
    $filename =~ s/.job$//g;
    print "Script handling $filename\n";
    print `perl $script $filename > $filename.jobout 2>&1`
}

sub eachFile {
    my $filename = $_;
    my $fullpath = $File::Find::name;
    #remember that File::Find changes your CWD, 
    #so you can call open with just $_

    return unless $filename =~ /\.job$/;
    
    unlink $filename;

    print "Got it\n";
    doRender($filename);	 
}

sub scan() {
    -d $root or die "Scan place is not a directory";
    print "Starting scan...\n";
    while ( 2 ) {
	sleep 2;
	find (\&eachFile, $root);
    }
}

scan();

__END__
=head1 NAME
sample - Using Getopt::Long and Pod::Usage
=head1 SYNOPSIS
sample [options] [file ...]
Options:
-help brief help message
-man full documentation
=head1 OPTIONS
=over 8
=item B<-help>
Print a brief help message and exits.
=item B<-man>
Prints the manual page and exits.
=back
=head1 DESCRIPTION
B<This program> will read the given input file(s) and do something
useful with the contents thereof.
=cut
