import Avatar from "@mui/material/Avatar";

function stringToColor(string: string) {
  let hash = 0;
  let i;

  /* eslint-disable no-bitwise */
  for (i = 0; i < string.length; i += 1) {
    hash = string.charCodeAt(i) + ((hash << 5) - hash);
  }

  let color = "#";

  for (i = 0; i < 3; i += 1) {
    const value = (hash >> (i * 8)) & 0xff;
    color += `00${value.toString(16)}`.slice(-2);
  }
  /* eslint-enable no-bitwise */

  return color;
}

function stringAvatar(name: string) {
  const parts = name.trim().split(" ");
  let initials = "";
  if (parts.length === 1) {
    initials = parts[0][0]?.toUpperCase() ?? "";
  } else {
    initials = `${parts[0][0]?.toUpperCase() ?? ""}${
      parts[1][0]?.toUpperCase() ?? ""
    }`;
  }
  return {
    sx: {
      bgcolor: stringToColor(name),
    },
    children: initials,
  };
}

export default function CustomAvatar({ name }: { name: string }) {
  const avatarProps = stringAvatar(name);
  return (
    <Avatar
      {...avatarProps}
      sx={{ ...avatarProps.sx, width: 28, height: 28, fontSize: 16 }}
    />
  );
}
