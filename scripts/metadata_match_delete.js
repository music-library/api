const fs = require('fs');
const path = require('path');

const metadata = require('./metadata.json');
const newMetadata = [];

//
// Match track and delete from metadata
//
console.log(`Match track and delete from metadata
`);

// Options
const DRY_RUN = false;
const MATCH_PATTERN = /(\[1958\] The Lady In Red)/i;

for (const key in metadata) {
	const track = metadata[key];
	const { path: trackPath, metadata: meta } = track;
	const { title, artist } = meta;

	// If path matches, skip
	if (trackPath.match(MATCH_PATTERN)) {
		console.log(`Matched: ${artist} - ${title}`);
		continue;
	}

	newMetadata.push(track);
}

console.log('');
console.log('Tracks before:', metadata.length);
console.log('Tracks after: ', newMetadata.length);
if (DRY_RUN) console.log('DRY RUN: No changes made');


// Overwrite metadata
if (!DRY_RUN) fs.writeFileSync(path.resolve(__dirname, 'metadata.json'), JSON.stringify(newMetadata));
