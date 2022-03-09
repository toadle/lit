mdfind kMDItemContentType == \"com.apple.application-*\" | tr \"\n\" \"\0\" | xargs -0 mdls -name \"kMDItemDisplayName\"
