// https://stackoverflow.com/a/13935718/4467556

#pragma once

#include <vector>
#include <string>

namespace digitalvenue::core::base64 {
    typedef unsigned char BYTE;

    std::string encode(BYTE const *buf, unsigned int bufLen);

    std::vector<BYTE> decode(std::string const &);
}
