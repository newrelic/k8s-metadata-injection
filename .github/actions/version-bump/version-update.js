function incrementMajorVersion(versionSegments) {
  return `${Number(versionSegments[0]) + 1}.0.0`;
}

function incrementMinorVersion(versionSegments) {
  return `${versionSegments[0]}.${Number(versionSegments[1]) + 1}.0`;
}

function incrementPatchVersion(versionSegments) {
  return `${versionSegments[0]}.${versionSegments[1]}.${Number(versionSegments[2]) + 1}`;
}

function getVersionDifferenceTypeAndIncrement(version, current, updated) {
  const currentSegments = current.split('.');
  const updatedSegments = updated.split('.');
  const versionSegments = version.split('.');

  if (currentSegments[0] !== updatedSegments[0]) {
    return incrementMajorVersion(versionSegments);
  } else if (currentSegments[1] !== updatedSegments[1]) {
    return incrementMinorVersion(versionSegments);
  } else if (currentSegments[2] !== updatedSegments[2]) {
    return incrementPatchVersion(versionSegments);
  } else {
    return "error";
  }
}

const [, , version, current, updated] = process.argv;
let result = getVersionDifferenceTypeAndIncrement(version, current, updated);
console.log(result)