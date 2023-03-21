const fs = require('fs');
const path = require('path');

const metadata = require('./metadata.json');
const newTracks = [];
const newTracksMap = {};

//
// Match track and delete from metadata
//
console.log(`Match track and delete from metadata
`);

// Options
const DRY_RUN = true;
const MATCH_PATTERN = /(\Be Mine Tonight)/i;

for (const key in metadata.tracks) {
	const track = metadata.tracks[key];
	const { metadata: meta } = track;
	const { artist, album_artist, album, title, year } = meta;

	// If path matches, skip
	if (`${year} ${album} ${album_artist} ${artist} ${title}`.match(MATCH_PATTERN)) {
		console.log(`Matched: ${artist} - ${title}`);
		continue;
	}

	newTracks.push(track);
}

// Re-build tracks_map
for (const key in metadata.tracks) {
	const track = metadata.tracks[key];
	newTracksMap[track.id] = Number(key);
}

console.log('');
console.log('Tracks before:', metadata.tracks.length);
console.log('Tracks after: ', newTracks.length);
if (DRY_RUN) console.log('DRY RUN: No changes made');

// Overwrite metadata
if (!DRY_RUN) fs.writeFileSync(path.resolve(__dirname, 'metadata.json'), JSON.stringify({ ...metadata, tracks: newTracks, tracks_map: newTracksMap }));
