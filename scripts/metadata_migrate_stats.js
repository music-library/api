const fs = require('fs');
const path = require('path');

const metadata = require('./metadata.json');
const metadataLegacy = require('./metadata.legacy.json');
const newTracks = [];

//
// Migrate track stats from legacy metadata.json to new metadata.json
//
console.log(`Migrate track stats from legacy metadata
`);

// Options
const DRY_RUN = true;

for (const key in metadata.tracks) {
	const track = { ...metadata.tracks[key] };
	const { metadata: meta } = track;
	const { artist, title } = meta;

	// Find legacy track
	const legacyTrack = metadataLegacy.find(t => t.metadata.album === track.metadata.album && t.metadata.album_artist === track.metadata.album_artist && t.metadata.artist === track.metadata.artist && t.metadata.title === track.metadata.title);
	if (legacyTrack) {
		console.log(`Matched: ${artist} - ${title}`);
		track.metadata.duration = legacyTrack.metadata.duration;
		track.stats = { ...legacyTrack.stats };
	}

	newTracks.push(track);
}

console.log('');
console.log('Tracks before:', metadata.tracks.length);
console.log('Tracks after: ', newTracks.length);
if (DRY_RUN) console.log('DRY RUN: No changes made');

// Overwrite metadata
if (!DRY_RUN) fs.writeFileSync(path.resolve(__dirname, 'metadata.json'), JSON.stringify({ ...metadata, tracks: newTracks }));
