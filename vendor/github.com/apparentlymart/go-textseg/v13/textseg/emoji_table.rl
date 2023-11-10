# The following Ragel file was autogenerated with unicode2ragel.rb 
# from: https://www.unicode.org/Public/13.0.0/ucd/emoji/emoji-data.txt
#
# It defines ["Extended_Pictographic"].
#
# To use this, make sure that your alphtype is set to byte,
# and that your input is in utf8.

%%{
    machine Emoji;
    
    Extended_Pictographic = 
        0xC2 0xA9               #E0.6   [1] (©️)       copyright
      | 0xC2 0xAE               #E0.6   [1] (®️)       registered
      | 0xE2 0x80 0xBC          #E0.6   [1] (‼️)       double exclamation mark
      | 0xE2 0x81 0x89          #E0.6   [1] (⁉️)       exclamation question ...
      | 0xE2 0x84 0xA2          #E0.6   [1] (™️)       trade mark
      | 0xE2 0x84 0xB9          #E0.6   [1] (ℹ️)       information
      | 0xE2 0x86 0x94..0x99    #E0.6   [6] (↔️..↙️)    left-right arrow..do...
      | 0xE2 0x86 0xA9..0xAA    #E0.6   [2] (↩️..↪️)    right arrow curving ...
      | 0xE2 0x8C 0x9A..0x9B    #E0.6   [2] (⌚..⌛)    watch..hourglass done
      | 0xE2 0x8C 0xA8          #E1.0   [1] (⌨️)       keyboard
      | 0xE2 0x8E 0x88          #E0.0   [1] (⎈)       HELM SYMBOL
      | 0xE2 0x8F 0x8F          #E1.0   [1] (⏏️)       eject button
      | 0xE2 0x8F 0xA9..0xAC    #E0.6   [4] (⏩..⏬)    fast-forward button..f...
      | 0xE2 0x8F 0xAD..0xAE    #E0.7   [2] (⏭️..⏮️)    next track button..l...
      | 0xE2 0x8F 0xAF          #E1.0   [1] (⏯️)       play or pause button
      | 0xE2 0x8F 0xB0          #E0.6   [1] (⏰)       alarm clock
      | 0xE2 0x8F 0xB1..0xB2    #E1.0   [2] (⏱️..⏲️)    stopwatch..timer clock
      | 0xE2 0x8F 0xB3          #E0.6   [1] (⏳)       hourglass not done
      | 0xE2 0x8F 0xB8..0xBA    #E0.7   [3] (⏸️..⏺️)    pause button..record...
      | 0xE2 0x93 0x82          #E0.6   [1] (Ⓜ️)       circled M
      | 0xE2 0x96 0xAA..0xAB    #E0.6   [2] (▪️..▫️)    black small square.....
      | 0xE2 0x96 0xB6          #E0.6   [1] (▶️)       play button
      | 0xE2 0x97 0x80          #E0.6   [1] (◀️)       reverse button
      | 0xE2 0x97 0xBB..0xBE    #E0.6   [4] (◻️..◾)    white medium square.....
      | 0xE2 0x98 0x80..0x81    #E0.6   [2] (☀️..☁️)    sun..cloud
      | 0xE2 0x98 0x82..0x83    #E0.7   [2] (☂️..☃️)    umbrella..snowman
      | 0xE2 0x98 0x84          #E1.0   [1] (☄️)       comet
      | 0xE2 0x98 0x85          #E0.0   [1] (★)       BLACK STAR
      | 0xE2 0x98 0x87..0x8D    #E0.0   [7] (☇..☍)    LIGHTNING..OPPOSITION
      | 0xE2 0x98 0x8E          #E0.6   [1] (☎️)       telephone
      | 0xE2 0x98 0x8F..0x90    #E0.0   [2] (☏..☐)    WHITE TELEPHONE..BALLO...
      | 0xE2 0x98 0x91          #E0.6   [1] (☑️)       check box with check
      | 0xE2 0x98 0x92          #E0.0   [1] (☒)       BALLOT BOX WITH X
      | 0xE2 0x98 0x94..0x95    #E0.6   [2] (☔..☕)    umbrella with rain dro...
      | 0xE2 0x98 0x96..0x97    #E0.0   [2] (☖..☗)    WHITE SHOGI PIECE..BLA...
      | 0xE2 0x98 0x98          #E1.0   [1] (☘️)       shamrock
      | 0xE2 0x98 0x99..0x9C    #E0.0   [4] (☙..☜)    REVERSED ROTATED FLORA...
      | 0xE2 0x98 0x9D          #E0.6   [1] (☝️)       index pointing up
      | 0xE2 0x98 0x9E..0x9F    #E0.0   [2] (☞..☟)    WHITE RIGHT POINTING I...
      | 0xE2 0x98 0xA0          #E1.0   [1] (☠️)       skull and crossbones
      | 0xE2 0x98 0xA1          #E0.0   [1] (☡)       CAUTION SIGN
      | 0xE2 0x98 0xA2..0xA3    #E1.0   [2] (☢️..☣️)    radioactive..biohazard
      | 0xE2 0x98 0xA4..0xA5    #E0.0   [2] (☤..☥)    CADUCEUS..ANKH
      | 0xE2 0x98 0xA6          #E1.0   [1] (☦️)       orthodox cross
      | 0xE2 0x98 0xA7..0xA9    #E0.0   [3] (☧..☩)    CHI RHO..CROSS OF JERU...
      | 0xE2 0x98 0xAA          #E0.7   [1] (☪️)       star and crescent
      | 0xE2 0x98 0xAB..0xAD    #E0.0   [3] (☫..☭)    FARSI SYMBOL..HAMMER A...
      | 0xE2 0x98 0xAE          #E1.0   [1] (☮️)       peace symbol
      | 0xE2 0x98 0xAF          #E0.7   [1] (☯️)       yin yang
      | 0xE2 0x98 0xB0..0xB7    #E0.0   [8] (☰..☷)    TRIGRAM FOR HEAVEN..TR...
      | 0xE2 0x98 0xB8..0xB9    #E0.7   [2] (☸️..☹️)    wheel of dharma..fro...
      | 0xE2 0x98 0xBA          #E0.6   [1] (☺️)       smiling face
      | 0xE2 0x98 0xBB..0xBF    #E0.0   [5] (☻..☿)    BLACK SMILING FACE..ME...
      | 0xE2 0x99 0x80          #E4.0   [1] (♀️)       female sign
      | 0xE2 0x99 0x81          #E0.0   [1] (♁)       EARTH
      | 0xE2 0x99 0x82          #E4.0   [1] (♂️)       male sign
      | 0xE2 0x99 0x83..0x87    #E0.0   [5] (♃..♇)    JUPITER..PLUTO
      | 0xE2 0x99 0x88..0x93    #E0.6  [12] (♈..♓)    Aries..Pisces
      | 0xE2 0x99 0x94..0x9E    #E0.0  [11] (♔..♞)    WHITE CHESS KING..BLAC...
      | 0xE2 0x99 0x9F          #E11.0  [1] (♟️)       chess pawn
      | 0xE2 0x99 0xA0          #E0.6   [1] (♠️)       spade suit
      | 0xE2 0x99 0xA1..0xA2    #E0.0   [2] (♡..♢)    WHITE HEART SUIT..WHIT...
      | 0xE2 0x99 0xA3          #E0.6   [1] (♣️)       club suit
      | 0xE2 0x99 0xA4          #E0.0   [1] (♤)       WHITE SPADE SUIT
      | 0xE2 0x99 0xA5..0xA6    #E0.6   [2] (♥️..♦️)    heart suit..diamond ...
      | 0xE2 0x99 0xA7          #E0.0   [1] (♧)       WHITE CLUB SUIT
      | 0xE2 0x99 0xA8          #E0.6   [1] (♨️)       hot springs
      | 0xE2 0x99 0xA9..0xBA    #E0.0  [18] (♩..♺)    QUARTER NOTE..RECYCLIN...
      | 0xE2 0x99 0xBB          #E0.6   [1] (♻️)       recycling symbol
      | 0xE2 0x99 0xBC..0xBD    #E0.0   [2] (♼..♽)    RECYCLED PAPER SYMBOL....
      | 0xE2 0x99 0xBE          #E11.0  [1] (♾️)       infinity
      | 0xE2 0x99 0xBF          #E0.6   [1] (♿)       wheelchair symbol
      | 0xE2 0x9A 0x80..0x85    #E0.0   [6] (⚀..⚅)    DIE FACE-1..DIE FACE-6
      | 0xE2 0x9A 0x90..0x91    #E0.0   [2] (⚐..⚑)    WHITE FLAG..BLACK FLAG
      | 0xE2 0x9A 0x92          #E1.0   [1] (⚒️)       hammer and pick
      | 0xE2 0x9A 0x93          #E0.6   [1] (⚓)       anchor
      | 0xE2 0x9A 0x94          #E1.0   [1] (⚔️)       crossed swords
      | 0xE2 0x9A 0x95          #E4.0   [1] (⚕️)       medical symbol
      | 0xE2 0x9A 0x96..0x97    #E1.0   [2] (⚖️..⚗️)    balance scale..alembic
      | 0xE2 0x9A 0x98          #E0.0   [1] (⚘)       FLOWER
      | 0xE2 0x9A 0x99          #E1.0   [1] (⚙️)       gear
      | 0xE2 0x9A 0x9A          #E0.0   [1] (⚚)       STAFF OF HERMES
      | 0xE2 0x9A 0x9B..0x9C    #E1.0   [2] (⚛️..⚜️)    atom symbol..fleur-d...
      | 0xE2 0x9A 0x9D..0x9F    #E0.0   [3] (⚝..⚟)    OUTLINED WHITE STAR..T...
      | 0xE2 0x9A 0xA0..0xA1    #E0.6   [2] (⚠️..⚡)    warning..high voltage
      | 0xE2 0x9A 0xA2..0xA6    #E0.0   [5] (⚢..⚦)    DOUBLED FEMALE SIGN..M...
      | 0xE2 0x9A 0xA7          #E13.0  [1] (⚧️)       transgender symbol
      | 0xE2 0x9A 0xA8..0xA9    #E0.0   [2] (⚨..⚩)    VERTICAL MALE WITH STR...
      | 0xE2 0x9A 0xAA..0xAB    #E0.6   [2] (⚪..⚫)    white circle..black ci...
      | 0xE2 0x9A 0xAC..0xAF    #E0.0   [4] (⚬..⚯)    MEDIUM SMALL WHITE CIR...
      | 0xE2 0x9A 0xB0..0xB1    #E1.0   [2] (⚰️..⚱️)    coffin..funeral urn
      | 0xE2 0x9A 0xB2..0xBC    #E0.0  [11] (⚲..⚼)    NEUTER..SESQUIQUADRATE
      | 0xE2 0x9A 0xBD..0xBE    #E0.6   [2] (⚽..⚾)    soccer ball..baseball
      | 0xE2 0x9A 0xBF..0xFF    #E0.0   [5] (⚿..⛃)    SQUARED KEY..BLACK DRA...
      | 0xE2 0x9B 0x00..0x83    #
      | 0xE2 0x9B 0x84..0x85    #E0.6   [2] (⛄..⛅)    snowman without snow.....
      | 0xE2 0x9B 0x86..0x87    #E0.0   [2] (⛆..⛇)    RAIN..BLACK SNOWMAN
      | 0xE2 0x9B 0x88          #E0.7   [1] (⛈️)       cloud with lightning ...
      | 0xE2 0x9B 0x89..0x8D    #E0.0   [5] (⛉..⛍)    TURNED WHITE SHOGI PIE...
      | 0xE2 0x9B 0x8E          #E0.6   [1] (⛎)       Ophiuchus
      | 0xE2 0x9B 0x8F          #E0.7   [1] (⛏️)       pick
      | 0xE2 0x9B 0x90          #E0.0   [1] (⛐)       CAR SLIDING
      | 0xE2 0x9B 0x91          #E0.7   [1] (⛑️)       rescue worker’s helmet
      | 0xE2 0x9B 0x92          #E0.0   [1] (⛒)       CIRCLED CROSSING LANES
      | 0xE2 0x9B 0x93          #E0.7   [1] (⛓️)       chains
      | 0xE2 0x9B 0x94          #E0.6   [1] (⛔)       no entry
      | 0xE2 0x9B 0x95..0xA8    #E0.0  [20] (⛕..⛨)    ALTERNATE ONE-WAY LEFT...
      | 0xE2 0x9B 0xA9          #E0.7   [1] (⛩️)       shinto shrine
      | 0xE2 0x9B 0xAA          #E0.6   [1] (⛪)       church
      | 0xE2 0x9B 0xAB..0xAF    #E0.0   [5] (⛫..⛯)    CASTLE..MAP SYMBOL FOR...
      | 0xE2 0x9B 0xB0..0xB1    #E0.7   [2] (⛰️..⛱️)    mountain..umbrella o...
      | 0xE2 0x9B 0xB2..0xB3    #E0.6   [2] (⛲..⛳)    fountain..flag in hole
      | 0xE2 0x9B 0xB4          #E0.7   [1] (⛴️)       ferry
      | 0xE2 0x9B 0xB5          #E0.6   [1] (⛵)       sailboat
      | 0xE2 0x9B 0xB6          #E0.0   [1] (⛶)       SQUARE FOUR CORNERS
      | 0xE2 0x9B 0xB7..0xB9    #E0.7   [3] (⛷️..⛹️)    skier..person bounci...
      | 0xE2 0x9B 0xBA          #E0.6   [1] (⛺)       tent
      | 0xE2 0x9B 0xBB..0xBC    #E0.0   [2] (⛻..⛼)    JAPANESE BANK SYMBOL.....
      | 0xE2 0x9B 0xBD          #E0.6   [1] (⛽)       fuel pump
      | 0xE2 0x9B 0xBE..0xFF    #E0.0   [4] (⛾..✁)    CUP ON BLACK SQUARE..U...
      | 0xE2 0x9C 0x00..0x81    #
      | 0xE2 0x9C 0x82          #E0.6   [1] (✂️)       scissors
      | 0xE2 0x9C 0x83..0x84    #E0.0   [2] (✃..✄)    LOWER BLADE SCISSORS.....
      | 0xE2 0x9C 0x85          #E0.6   [1] (✅)       check mark button
      | 0xE2 0x9C 0x88..0x8C    #E0.6   [5] (✈️..✌️)    airplane..victory hand
      | 0xE2 0x9C 0x8D          #E0.7   [1] (✍️)       writing hand
      | 0xE2 0x9C 0x8E          #E0.0   [1] (✎)       LOWER RIGHT PENCIL
      | 0xE2 0x9C 0x8F          #E0.6   [1] (✏️)       pencil
      | 0xE2 0x9C 0x90..0x91    #E0.0   [2] (✐..✑)    UPPER RIGHT PENCIL..WH...
      | 0xE2 0x9C 0x92          #E0.6   [1] (✒️)       black nib
      | 0xE2 0x9C 0x94          #E0.6   [1] (✔️)       check mark
      | 0xE2 0x9C 0x96          #E0.6   [1] (✖️)       multiply
      | 0xE2 0x9C 0x9D          #E0.7   [1] (✝️)       latin cross
      | 0xE2 0x9C 0xA1          #E0.7   [1] (✡️)       star of David
      | 0xE2 0x9C 0xA8          #E0.6   [1] (✨)       sparkles
      | 0xE2 0x9C 0xB3..0xB4    #E0.6   [2] (✳️..✴️)    eight-spoked asteris...
      | 0xE2 0x9D 0x84          #E0.6   [1] (❄️)       snowflake
      | 0xE2 0x9D 0x87          #E0.6   [1] (❇️)       sparkle
      | 0xE2 0x9D 0x8C          #E0.6   [1] (❌)       cross mark
      | 0xE2 0x9D 0x8E          #E0.6   [1] (❎)       cross mark button
      | 0xE2 0x9D 0x93..0x95    #E0.6   [3] (❓..❕)    question mark..white e...
      | 0xE2 0x9D 0x97          #E0.6   [1] (❗)       exclamation mark
      | 0xE2 0x9D 0xA3          #E1.0   [1] (❣️)       heart exclamation
      | 0xE2 0x9D 0xA4          #E0.6   [1] (❤️)       red heart
      | 0xE2 0x9D 0xA5..0xA7    #E0.0   [3] (❥..❧)    ROTATED HEAVY BLACK HE...
      | 0xE2 0x9E 0x95..0x97    #E0.6   [3] (➕..➗)    plus..divide
      | 0xE2 0x9E 0xA1          #E0.6   [1] (➡️)       right arrow
      | 0xE2 0x9E 0xB0          #E0.6   [1] (➰)       curly loop
      | 0xE2 0x9E 0xBF          #E1.0   [1] (➿)       double curly loop
      | 0xE2 0xA4 0xB4..0xB5    #E0.6   [2] (⤴️..⤵️)    right arrow curving ...
      | 0xE2 0xAC 0x85..0x87    #E0.6   [3] (⬅️..⬇️)    left arrow..down arrow
      | 0xE2 0xAC 0x9B..0x9C    #E0.6   [2] (⬛..⬜)    black large square..wh...
      | 0xE2 0xAD 0x90          #E0.6   [1] (⭐)       star
      | 0xE2 0xAD 0x95          #E0.6   [1] (⭕)       hollow red circle
      | 0xE3 0x80 0xB0          #E0.6   [1] (〰️)       wavy dash
      | 0xE3 0x80 0xBD          #E0.6   [1] (〽️)       part alternation mark
      | 0xE3 0x8A 0x97          #E0.6   [1] (㊗️)       Japanese “congratulat...
      | 0xE3 0x8A 0x99          #E0.6   [1] (㊙️)       Japanese “secret” button
      | 0xF0 0x9F 0x80 0x80..0x83  #E0.0   [4] (🀀..🀃)    MAHJONG TILE EAST W...
      | 0xF0 0x9F 0x80 0x84     #E0.6   [1] (🀄)       mahjong red dragon
      | 0xF0 0x9F 0x80 0x85..0xFF        #E0.0 [202] (🀅..🃎)    MAHJONG TILE ...
      | 0xF0 0x9F 0x81..0x82 0x00..0xFF  #
      | 0xF0 0x9F 0x83 0x00..0x8E        #
      | 0xF0 0x9F 0x83 0x8F     #E0.6   [1] (🃏)       joker
      | 0xF0 0x9F 0x83 0x90..0xBF  #E0.0  [48] (🃐..🃿)    <reserved-1F0D0>..<...
      | 0xF0 0x9F 0x84 0x8D..0x8F  #E0.0   [3] (🄍..🄏)    CIRCLED ZERO WITH S...
      | 0xF0 0x9F 0x84 0xAF     #E0.0   [1] (🄯)       COPYLEFT SYMBOL
      | 0xF0 0x9F 0x85 0xAC..0xAF  #E0.0   [4] (🅬..🅯)    RAISED MR SIGN..CIR...
      | 0xF0 0x9F 0x85 0xB0..0xB1  #E0.6   [2] (🅰️..🅱️)    A button (blood t...
      | 0xF0 0x9F 0x85 0xBE..0xBF  #E0.6   [2] (🅾️..🅿️)    O button (blood t...
      | 0xF0 0x9F 0x86 0x8E     #E0.6   [1] (🆎)       AB button (blood type)
      | 0xF0 0x9F 0x86 0x91..0x9A  #E0.6  [10] (🆑..🆚)    CL button..VS button
      | 0xF0 0x9F 0x86 0xAD..0xFF  #E0.0  [57] (🆭..🇥)    MASK WORK SYMBOL..<...
      | 0xF0 0x9F 0x87 0x00..0xA5  #
      | 0xF0 0x9F 0x88 0x81..0x82  #E0.6   [2] (🈁..🈂️)    Japanese “here” bu...
      | 0xF0 0x9F 0x88 0x83..0x8F  #E0.0  [13] (🈃..🈏)    <reserved-1F203>..<...
      | 0xF0 0x9F 0x88 0x9A     #E0.6   [1] (🈚)       Japanese “free of char...
      | 0xF0 0x9F 0x88 0xAF     #E0.6   [1] (🈯)       Japanese “reserved” bu...
      | 0xF0 0x9F 0x88 0xB2..0xBA  #E0.6   [9] (🈲..🈺)    Japanese “prohibite...
      | 0xF0 0x9F 0x88 0xBC..0xBF  #E0.0   [4] (🈼..🈿)    <reserved-1F23C>..<...
      | 0xF0 0x9F 0x89 0x89..0x8F  #E0.0   [7] (🉉..🉏)    <reserved-1F249>..<...
      | 0xF0 0x9F 0x89 0x90..0x91  #E0.6   [2] (🉐..🉑)    Japanese “bargain” ...
      | 0xF0 0x9F 0x89 0x92..0xFF        #E0.0 [174] (🉒..🋿)    <reserved-1F2...
      | 0xF0 0x9F 0x8A..0x8A 0x00..0xFF  #
      | 0xF0 0x9F 0x8B 0x00..0xBF        #
      | 0xF0 0x9F 0x8C 0x80..0x8C  #E0.6  [13] (🌀..🌌)    cyclone..milky way
      | 0xF0 0x9F 0x8C 0x8D..0x8E  #E0.7   [2] (🌍..🌎)    globe showing Europ...
      | 0xF0 0x9F 0x8C 0x8F     #E0.6   [1] (🌏)       globe showing Asia-Aus...
      | 0xF0 0x9F 0x8C 0x90     #E1.0   [1] (🌐)       globe with meridians
      | 0xF0 0x9F 0x8C 0x91     #E0.6   [1] (🌑)       new moon
      | 0xF0 0x9F 0x8C 0x92     #E1.0   [1] (🌒)       waxing crescent moon
      | 0xF0 0x9F 0x8C 0x93..0x95  #E0.6   [3] (🌓..🌕)    first quarter moon....
      | 0xF0 0x9F 0x8C 0x96..0x98  #E1.0   [3] (🌖..🌘)    waning gibbous moon...
      | 0xF0 0x9F 0x8C 0x99     #E0.6   [1] (🌙)       crescent moon
      | 0xF0 0x9F 0x8C 0x9A     #E1.0   [1] (🌚)       new moon face
      | 0xF0 0x9F 0x8C 0x9B     #E0.6   [1] (🌛)       first quarter moon face
      | 0xF0 0x9F 0x8C 0x9C     #E0.7   [1] (🌜)       last quarter moon face
      | 0xF0 0x9F 0x8C 0x9D..0x9E  #E1.0   [2] (🌝..🌞)    full moon face..sun...
      | 0xF0 0x9F 0x8C 0x9F..0xA0  #E0.6   [2] (🌟..🌠)    glowing star..shoot...
      | 0xF0 0x9F 0x8C 0xA1     #E0.7   [1] (🌡️)       thermometer
      | 0xF0 0x9F 0x8C 0xA2..0xA3  #E0.0   [2] (🌢..🌣)    BLACK DROPLET..WHIT...
      | 0xF0 0x9F 0x8C 0xA4..0xAC  #E0.7   [9] (🌤️..🌬️)    sun behind small ...
      | 0xF0 0x9F 0x8C 0xAD..0xAF  #E1.0   [3] (🌭..🌯)    hot dog..burrito
      | 0xF0 0x9F 0x8C 0xB0..0xB1  #E0.6   [2] (🌰..🌱)    chestnut..seedling
      | 0xF0 0x9F 0x8C 0xB2..0xB3  #E1.0   [2] (🌲..🌳)    evergreen tree..dec...
      | 0xF0 0x9F 0x8C 0xB4..0xB5  #E0.6   [2] (🌴..🌵)    palm tree..cactus
      | 0xF0 0x9F 0x8C 0xB6     #E0.7   [1] (🌶️)       hot pepper
      | 0xF0 0x9F 0x8C 0xB7..0xFF  #E0.6  [20] (🌷..🍊)    tulip..tangerine
      | 0xF0 0x9F 0x8D 0x00..0x8A  #
      | 0xF0 0x9F 0x8D 0x8B     #E1.0   [1] (🍋)       lemon
      | 0xF0 0x9F 0x8D 0x8C..0x8F  #E0.6   [4] (🍌..🍏)    banana..green apple
      | 0xF0 0x9F 0x8D 0x90     #E1.0   [1] (🍐)       pear
      | 0xF0 0x9F 0x8D 0x91..0xBB  #E0.6  [43] (🍑..🍻)    peach..clinking bee...
      | 0xF0 0x9F 0x8D 0xBC     #E1.0   [1] (🍼)       baby bottle
      | 0xF0 0x9F 0x8D 0xBD     #E0.7   [1] (🍽️)       fork and knife with p...
      | 0xF0 0x9F 0x8D 0xBE..0xBF  #E1.0   [2] (🍾..🍿)    bottle with popping...
      | 0xF0 0x9F 0x8E 0x80..0x93  #E0.6  [20] (🎀..🎓)    ribbon..graduation cap
      | 0xF0 0x9F 0x8E 0x94..0x95  #E0.0   [2] (🎔..🎕)    HEART WITH TIP ON T...
      | 0xF0 0x9F 0x8E 0x96..0x97  #E0.7   [2] (🎖️..🎗️)    military medal..r...
      | 0xF0 0x9F 0x8E 0x98     #E0.0   [1] (🎘)       MUSICAL KEYBOARD WITH ...
      | 0xF0 0x9F 0x8E 0x99..0x9B  #E0.7   [3] (🎙️..🎛️)    studio microphone...
      | 0xF0 0x9F 0x8E 0x9C..0x9D  #E0.0   [2] (🎜..🎝)    BEAMED ASCENDING MU...
      | 0xF0 0x9F 0x8E 0x9E..0x9F  #E0.7   [2] (🎞️..🎟️)    film frames..admi...
      | 0xF0 0x9F 0x8E 0xA0..0xFF  #E0.6  [37] (🎠..🏄)    carousel horse..per...
      | 0xF0 0x9F 0x8F 0x00..0x84  #
      | 0xF0 0x9F 0x8F 0x85     #E1.0   [1] (🏅)       sports medal
      | 0xF0 0x9F 0x8F 0x86     #E0.6   [1] (🏆)       trophy
      | 0xF0 0x9F 0x8F 0x87     #E1.0   [1] (🏇)       horse racing
      | 0xF0 0x9F 0x8F 0x88     #E0.6   [1] (🏈)       american football
      | 0xF0 0x9F 0x8F 0x89     #E1.0   [1] (🏉)       rugby football
      | 0xF0 0x9F 0x8F 0x8A     #E0.6   [1] (🏊)       person swimming
      | 0xF0 0x9F 0x8F 0x8B..0x8E  #E0.7   [4] (🏋️..🏎️)    person lifting we...
      | 0xF0 0x9F 0x8F 0x8F..0x93  #E1.0   [5] (🏏..🏓)    cricket game..ping ...
      | 0xF0 0x9F 0x8F 0x94..0x9F  #E0.7  [12] (🏔️..🏟️)    snow-capped mount...
      | 0xF0 0x9F 0x8F 0xA0..0xA3  #E0.6   [4] (🏠..🏣)    house..Japanese pos...
      | 0xF0 0x9F 0x8F 0xA4     #E1.0   [1] (🏤)       post office
      | 0xF0 0x9F 0x8F 0xA5..0xB0  #E0.6  [12] (🏥..🏰)    hospital..castle
      | 0xF0 0x9F 0x8F 0xB1..0xB2  #E0.0   [2] (🏱..🏲)    WHITE PENNANT..BLAC...
      | 0xF0 0x9F 0x8F 0xB3     #E0.7   [1] (🏳️)       white flag
      | 0xF0 0x9F 0x8F 0xB4     #E1.0   [1] (🏴)       black flag
      | 0xF0 0x9F 0x8F 0xB5     #E0.7   [1] (🏵️)       rosette
      | 0xF0 0x9F 0x8F 0xB6     #E0.0   [1] (🏶)       BLACK ROSETTE
      | 0xF0 0x9F 0x8F 0xB7     #E0.7   [1] (🏷️)       label
      | 0xF0 0x9F 0x8F 0xB8..0xBA  #E1.0   [3] (🏸..🏺)    badminton..amphora
      | 0xF0 0x9F 0x90 0x80..0x87  #E1.0   [8] (🐀..🐇)    rat..rabbit
      | 0xF0 0x9F 0x90 0x88     #E0.7   [1] (🐈)       cat
      | 0xF0 0x9F 0x90 0x89..0x8B  #E1.0   [3] (🐉..🐋)    dragon..whale
      | 0xF0 0x9F 0x90 0x8C..0x8E  #E0.6   [3] (🐌..🐎)    snail..horse
      | 0xF0 0x9F 0x90 0x8F..0x90  #E1.0   [2] (🐏..🐐)    ram..goat
      | 0xF0 0x9F 0x90 0x91..0x92  #E0.6   [2] (🐑..🐒)    ewe..monkey
      | 0xF0 0x9F 0x90 0x93     #E1.0   [1] (🐓)       rooster
      | 0xF0 0x9F 0x90 0x94     #E0.6   [1] (🐔)       chicken
      | 0xF0 0x9F 0x90 0x95     #E0.7   [1] (🐕)       dog
      | 0xF0 0x9F 0x90 0x96     #E1.0   [1] (🐖)       pig
      | 0xF0 0x9F 0x90 0x97..0xA9  #E0.6  [19] (🐗..🐩)    boar..poodle
      | 0xF0 0x9F 0x90 0xAA     #E1.0   [1] (🐪)       camel
      | 0xF0 0x9F 0x90 0xAB..0xBE  #E0.6  [20] (🐫..🐾)    two-hump camel..paw...
      | 0xF0 0x9F 0x90 0xBF     #E0.7   [1] (🐿️)       chipmunk
      | 0xF0 0x9F 0x91 0x80     #E0.6   [1] (👀)       eyes
      | 0xF0 0x9F 0x91 0x81     #E0.7   [1] (👁️)       eye
      | 0xF0 0x9F 0x91 0x82..0xA4  #E0.6  [35] (👂..👤)    ear..bust in silhou...
      | 0xF0 0x9F 0x91 0xA5     #E1.0   [1] (👥)       busts in silhouette
      | 0xF0 0x9F 0x91 0xA6..0xAB  #E0.6   [6] (👦..👫)    boy..woman and man ...
      | 0xF0 0x9F 0x91 0xAC..0xAD  #E1.0   [2] (👬..👭)    men holding hands.....
      | 0xF0 0x9F 0x91 0xAE..0xFF  #E0.6  [63] (👮..💬)    police officer..spe...
      | 0xF0 0x9F 0x92 0x00..0xAC  #
      | 0xF0 0x9F 0x92 0xAD     #E1.0   [1] (💭)       thought balloon
      | 0xF0 0x9F 0x92 0xAE..0xB5  #E0.6   [8] (💮..💵)    white flower..dolla...
      | 0xF0 0x9F 0x92 0xB6..0xB7  #E1.0   [2] (💶..💷)    euro banknote..poun...
      | 0xF0 0x9F 0x92 0xB8..0xFF  #E0.6  [52] (💸..📫)    money with wings..c...
      | 0xF0 0x9F 0x93 0x00..0xAB  #
      | 0xF0 0x9F 0x93 0xAC..0xAD  #E0.7   [2] (📬..📭)    open mailbox with r...
      | 0xF0 0x9F 0x93 0xAE     #E0.6   [1] (📮)       postbox
      | 0xF0 0x9F 0x93 0xAF     #E1.0   [1] (📯)       postal horn
      | 0xF0 0x9F 0x93 0xB0..0xB4  #E0.6   [5] (📰..📴)    newspaper..mobile p...
      | 0xF0 0x9F 0x93 0xB5     #E1.0   [1] (📵)       no mobile phones
      | 0xF0 0x9F 0x93 0xB6..0xB7  #E0.6   [2] (📶..📷)    antenna bars..camera
      | 0xF0 0x9F 0x93 0xB8     #E1.0   [1] (📸)       camera with flash
      | 0xF0 0x9F 0x93 0xB9..0xBC  #E0.6   [4] (📹..📼)    video camera..video...
      | 0xF0 0x9F 0x93 0xBD     #E0.7   [1] (📽️)       film projector
      | 0xF0 0x9F 0x93 0xBE     #E0.0   [1] (📾)       PORTABLE STEREO
      | 0xF0 0x9F 0x93 0xBF..0xFF  #E1.0   [4] (📿..🔂)    prayer beads..repea...
      | 0xF0 0x9F 0x94 0x00..0x82  #
      | 0xF0 0x9F 0x94 0x83     #E0.6   [1] (🔃)       clockwise vertical arrows
      | 0xF0 0x9F 0x94 0x84..0x87  #E1.0   [4] (🔄..🔇)    counterclockwise ar...
      | 0xF0 0x9F 0x94 0x88     #E0.7   [1] (🔈)       speaker low volume
      | 0xF0 0x9F 0x94 0x89     #E1.0   [1] (🔉)       speaker medium volume
      | 0xF0 0x9F 0x94 0x8A..0x94  #E0.6  [11] (🔊..🔔)    speaker high volume...
      | 0xF0 0x9F 0x94 0x95     #E1.0   [1] (🔕)       bell with slash
      | 0xF0 0x9F 0x94 0x96..0xAB  #E0.6  [22] (🔖..🔫)    bookmark..pistol
      | 0xF0 0x9F 0x94 0xAC..0xAD  #E1.0   [2] (🔬..🔭)    microscope..telescope
      | 0xF0 0x9F 0x94 0xAE..0xBD  #E0.6  [16] (🔮..🔽)    crystal ball..downw...
      | 0xF0 0x9F 0x95 0x86..0x88  #E0.0   [3] (🕆..🕈)    WHITE LATIN CROSS.....
      | 0xF0 0x9F 0x95 0x89..0x8A  #E0.7   [2] (🕉️..🕊️)    om..dove
      | 0xF0 0x9F 0x95 0x8B..0x8E  #E1.0   [4] (🕋..🕎)    kaaba..menorah
      | 0xF0 0x9F 0x95 0x8F     #E0.0   [1] (🕏)       BOWL OF HYGIEIA
      | 0xF0 0x9F 0x95 0x90..0x9B  #E0.6  [12] (🕐..🕛)    one o’clock..twelve...
      | 0xF0 0x9F 0x95 0x9C..0xA7  #E0.7  [12] (🕜..🕧)    one-thirty..twelve-...
      | 0xF0 0x9F 0x95 0xA8..0xAE  #E0.0   [7] (🕨..🕮)    RIGHT SPEAKER..BOOK
      | 0xF0 0x9F 0x95 0xAF..0xB0  #E0.7   [2] (🕯️..🕰️)    candle..mantelpie...
      | 0xF0 0x9F 0x95 0xB1..0xB2  #E0.0   [2] (🕱..🕲)    BLACK SKULL AND CRO...
      | 0xF0 0x9F 0x95 0xB3..0xB9  #E0.7   [7] (🕳️..🕹️)    hole..joystick
      | 0xF0 0x9F 0x95 0xBA     #E3.0   [1] (🕺)       man dancing
      | 0xF0 0x9F 0x95 0xBB..0xFF  #E0.0  [12] (🕻..🖆)    LEFT HAND TELEPHONE...
      | 0xF0 0x9F 0x96 0x00..0x86  #
      | 0xF0 0x9F 0x96 0x87     #E0.7   [1] (🖇️)       linked paperclips
      | 0xF0 0x9F 0x96 0x88..0x89  #E0.0   [2] (🖈..🖉)    BLACK PUSHPIN..LOWE...
      | 0xF0 0x9F 0x96 0x8A..0x8D  #E0.7   [4] (🖊️..🖍️)    pen..crayon
      | 0xF0 0x9F 0x96 0x8E..0x8F  #E0.0   [2] (🖎..🖏)    LEFT WRITING HAND.....
      | 0xF0 0x9F 0x96 0x90     #E0.7   [1] (🖐️)       hand with fingers spl...
      | 0xF0 0x9F 0x96 0x91..0x94  #E0.0   [4] (🖑..🖔)    REVERSED RAISED HAN...
      | 0xF0 0x9F 0x96 0x95..0x96  #E1.0   [2] (🖕..🖖)    middle finger..vulc...
      | 0xF0 0x9F 0x96 0x97..0xA3  #E0.0  [13] (🖗..🖣)    WHITE DOWN POINTING...
      | 0xF0 0x9F 0x96 0xA4     #E3.0   [1] (🖤)       black heart
      | 0xF0 0x9F 0x96 0xA5     #E0.7   [1] (🖥️)       desktop computer
      | 0xF0 0x9F 0x96 0xA6..0xA7  #E0.0   [2] (🖦..🖧)    KEYBOARD AND MOUSE....
      | 0xF0 0x9F 0x96 0xA8     #E0.7   [1] (🖨️)       printer
      | 0xF0 0x9F 0x96 0xA9..0xB0  #E0.0   [8] (🖩..🖰)    POCKET CALCULATOR.....
      | 0xF0 0x9F 0x96 0xB1..0xB2  #E0.7   [2] (🖱️..🖲️)    computer mouse..t...
      | 0xF0 0x9F 0x96 0xB3..0xBB  #E0.0   [9] (🖳..🖻)    OLD PERSONAL COMPUT...
      | 0xF0 0x9F 0x96 0xBC     #E0.7   [1] (🖼️)       framed picture
      | 0xF0 0x9F 0x96 0xBD..0xFF  #E0.0   [5] (🖽..🗁)    FRAME WITH TILES..O...
      | 0xF0 0x9F 0x97 0x00..0x81  #
      | 0xF0 0x9F 0x97 0x82..0x84  #E0.7   [3] (🗂️..🗄️)    card index divide...
      | 0xF0 0x9F 0x97 0x85..0x90  #E0.0  [12] (🗅..🗐)    EMPTY NOTE..PAGES
      | 0xF0 0x9F 0x97 0x91..0x93  #E0.7   [3] (🗑️..🗓️)    wastebasket..spir...
      | 0xF0 0x9F 0x97 0x94..0x9B  #E0.0   [8] (🗔..🗛)    DESKTOP WINDOW..DEC...
      | 0xF0 0x9F 0x97 0x9C..0x9E  #E0.7   [3] (🗜️..🗞️)    clamp..rolled-up ...
      | 0xF0 0x9F 0x97 0x9F..0xA0  #E0.0   [2] (🗟..🗠)    PAGE WITH CIRCLED T...
      | 0xF0 0x9F 0x97 0xA1     #E0.7   [1] (🗡️)       dagger
      | 0xF0 0x9F 0x97 0xA2     #E0.0   [1] (🗢)       LIPS
      | 0xF0 0x9F 0x97 0xA3     #E0.7   [1] (🗣️)       speaking head
      | 0xF0 0x9F 0x97 0xA4..0xA7  #E0.0   [4] (🗤..🗧)    THREE RAYS ABOVE..T...
      | 0xF0 0x9F 0x97 0xA8     #E2.0   [1] (🗨️)       left speech bubble
      | 0xF0 0x9F 0x97 0xA9..0xAE  #E0.0   [6] (🗩..🗮)    RIGHT SPEECH BUBBLE...
      | 0xF0 0x9F 0x97 0xAF     #E0.7   [1] (🗯️)       right anger bubble
      | 0xF0 0x9F 0x97 0xB0..0xB2  #E0.0   [3] (🗰..🗲)    MOOD BUBBLE..LIGHTN...
      | 0xF0 0x9F 0x97 0xB3     #E0.7   [1] (🗳️)       ballot box with ballot
      | 0xF0 0x9F 0x97 0xB4..0xB9  #E0.0   [6] (🗴..🗹)    BALLOT SCRIPT X..BA...
      | 0xF0 0x9F 0x97 0xBA     #E0.7   [1] (🗺️)       world map
      | 0xF0 0x9F 0x97 0xBB..0xBF  #E0.6   [5] (🗻..🗿)    mount fuji..moai
      | 0xF0 0x9F 0x98 0x80     #E1.0   [1] (😀)       grinning face
      | 0xF0 0x9F 0x98 0x81..0x86  #E0.6   [6] (😁..😆)    beaming face with s...
      | 0xF0 0x9F 0x98 0x87..0x88  #E1.0   [2] (😇..😈)    smiling face with h...
      | 0xF0 0x9F 0x98 0x89..0x8D  #E0.6   [5] (😉..😍)    winking face..smili...
      | 0xF0 0x9F 0x98 0x8E     #E1.0   [1] (😎)       smiling face with sung...
      | 0xF0 0x9F 0x98 0x8F     #E0.6   [1] (😏)       smirking face
      | 0xF0 0x9F 0x98 0x90     #E0.7   [1] (😐)       neutral face
      | 0xF0 0x9F 0x98 0x91     #E1.0   [1] (😑)       expressionless face
      | 0xF0 0x9F 0x98 0x92..0x94  #E0.6   [3] (😒..😔)    unamused face..pens...
      | 0xF0 0x9F 0x98 0x95     #E1.0   [1] (😕)       confused face
      | 0xF0 0x9F 0x98 0x96     #E0.6   [1] (😖)       confounded face
      | 0xF0 0x9F 0x98 0x97     #E1.0   [1] (😗)       kissing face
      | 0xF0 0x9F 0x98 0x98     #E0.6   [1] (😘)       face blowing a kiss
      | 0xF0 0x9F 0x98 0x99     #E1.0   [1] (😙)       kissing face with smil...
      | 0xF0 0x9F 0x98 0x9A     #E0.6   [1] (😚)       kissing face with clos...
      | 0xF0 0x9F 0x98 0x9B     #E1.0   [1] (😛)       face with tongue
      | 0xF0 0x9F 0x98 0x9C..0x9E  #E0.6   [3] (😜..😞)    winking face with t...
      | 0xF0 0x9F 0x98 0x9F     #E1.0   [1] (😟)       worried face
      | 0xF0 0x9F 0x98 0xA0..0xA5  #E0.6   [6] (😠..😥)    angry face..sad but...
      | 0xF0 0x9F 0x98 0xA6..0xA7  #E1.0   [2] (😦..😧)    frowning face with ...
      | 0xF0 0x9F 0x98 0xA8..0xAB  #E0.6   [4] (😨..😫)    fearful face..tired...
      | 0xF0 0x9F 0x98 0xAC     #E1.0   [1] (😬)       grimacing face
      | 0xF0 0x9F 0x98 0xAD     #E0.6   [1] (😭)       loudly crying face
      | 0xF0 0x9F 0x98 0xAE..0xAF  #E1.0   [2] (😮..😯)    face with open mout...
      | 0xF0 0x9F 0x98 0xB0..0xB3  #E0.6   [4] (😰..😳)    anxious face with s...
      | 0xF0 0x9F 0x98 0xB4     #E1.0   [1] (😴)       sleeping face
      | 0xF0 0x9F 0x98 0xB5     #E0.6   [1] (😵)       dizzy face
      | 0xF0 0x9F 0x98 0xB6     #E1.0   [1] (😶)       face without mouth
      | 0xF0 0x9F 0x98 0xB7..0xFF  #E0.6  [10] (😷..🙀)    face with medical m...
      | 0xF0 0x9F 0x99 0x00..0x80  #
      | 0xF0 0x9F 0x99 0x81..0x84  #E1.0   [4] (🙁..🙄)    slightly frowning f...
      | 0xF0 0x9F 0x99 0x85..0x8F  #E0.6  [11] (🙅..🙏)    person gesturing NO...
      | 0xF0 0x9F 0x9A 0x80     #E0.6   [1] (🚀)       rocket
      | 0xF0 0x9F 0x9A 0x81..0x82  #E1.0   [2] (🚁..🚂)    helicopter..locomotive
      | 0xF0 0x9F 0x9A 0x83..0x85  #E0.6   [3] (🚃..🚅)    railway car..bullet...
      | 0xF0 0x9F 0x9A 0x86     #E1.0   [1] (🚆)       train
      | 0xF0 0x9F 0x9A 0x87     #E0.6   [1] (🚇)       metro
      | 0xF0 0x9F 0x9A 0x88     #E1.0   [1] (🚈)       light rail
      | 0xF0 0x9F 0x9A 0x89     #E0.6   [1] (🚉)       station
      | 0xF0 0x9F 0x9A 0x8A..0x8B  #E1.0   [2] (🚊..🚋)    tram..tram car
      | 0xF0 0x9F 0x9A 0x8C     #E0.6   [1] (🚌)       bus
      | 0xF0 0x9F 0x9A 0x8D     #E0.7   [1] (🚍)       oncoming bus
      | 0xF0 0x9F 0x9A 0x8E     #E1.0   [1] (🚎)       trolleybus
      | 0xF0 0x9F 0x9A 0x8F     #E0.6   [1] (🚏)       bus stop
      | 0xF0 0x9F 0x9A 0x90     #E1.0   [1] (🚐)       minibus
      | 0xF0 0x9F 0x9A 0x91..0x93  #E0.6   [3] (🚑..🚓)    ambulance..police car
      | 0xF0 0x9F 0x9A 0x94     #E0.7   [1] (🚔)       oncoming police car
      | 0xF0 0x9F 0x9A 0x95     #E0.6   [1] (🚕)       taxi
      | 0xF0 0x9F 0x9A 0x96     #E1.0   [1] (🚖)       oncoming taxi
      | 0xF0 0x9F 0x9A 0x97     #E0.6   [1] (🚗)       automobile
      | 0xF0 0x9F 0x9A 0x98     #E0.7   [1] (🚘)       oncoming automobile
      | 0xF0 0x9F 0x9A 0x99..0x9A  #E0.6   [2] (🚙..🚚)    sport utility vehic...
      | 0xF0 0x9F 0x9A 0x9B..0xA1  #E1.0   [7] (🚛..🚡)    articulated lorry.....
      | 0xF0 0x9F 0x9A 0xA2     #E0.6   [1] (🚢)       ship
      | 0xF0 0x9F 0x9A 0xA3     #E1.0   [1] (🚣)       person rowing boat
      | 0xF0 0x9F 0x9A 0xA4..0xA5  #E0.6   [2] (🚤..🚥)    speedboat..horizont...
      | 0xF0 0x9F 0x9A 0xA6     #E1.0   [1] (🚦)       vertical traffic light
      | 0xF0 0x9F 0x9A 0xA7..0xAD  #E0.6   [7] (🚧..🚭)    construction..no sm...
      | 0xF0 0x9F 0x9A 0xAE..0xB1  #E1.0   [4] (🚮..🚱)    litter in bin sign....
      | 0xF0 0x9F 0x9A 0xB2     #E0.6   [1] (🚲)       bicycle
      | 0xF0 0x9F 0x9A 0xB3..0xB5  #E1.0   [3] (🚳..🚵)    no bicycles..person...
      | 0xF0 0x9F 0x9A 0xB6     #E0.6   [1] (🚶)       person walking
      | 0xF0 0x9F 0x9A 0xB7..0xB8  #E1.0   [2] (🚷..🚸)    no pedestrians..chi...
      | 0xF0 0x9F 0x9A 0xB9..0xBE  #E0.6   [6] (🚹..🚾)    men’s room..water c...
      | 0xF0 0x9F 0x9A 0xBF     #E1.0   [1] (🚿)       shower
      | 0xF0 0x9F 0x9B 0x80     #E0.6   [1] (🛀)       person taking bath
      | 0xF0 0x9F 0x9B 0x81..0x85  #E1.0   [5] (🛁..🛅)    bathtub..left luggage
      | 0xF0 0x9F 0x9B 0x86..0x8A  #E0.0   [5] (🛆..🛊)    TRIANGLE WITH ROUND...
      | 0xF0 0x9F 0x9B 0x8B     #E0.7   [1] (🛋️)       couch and lamp
      | 0xF0 0x9F 0x9B 0x8C     #E1.0   [1] (🛌)       person in bed
      | 0xF0 0x9F 0x9B 0x8D..0x8F  #E0.7   [3] (🛍️..🛏️)    shopping bags..bed
      | 0xF0 0x9F 0x9B 0x90     #E1.0   [1] (🛐)       place of worship
      | 0xF0 0x9F 0x9B 0x91..0x92  #E3.0   [2] (🛑..🛒)    stop sign..shopping...
      | 0xF0 0x9F 0x9B 0x93..0x94  #E0.0   [2] (🛓..🛔)    STUPA..PAGODA
      | 0xF0 0x9F 0x9B 0x95     #E12.0  [1] (🛕)       hindu temple
      | 0xF0 0x9F 0x9B 0x96..0x97  #E13.0  [2] (🛖..🛗)    hut..elevator
      | 0xF0 0x9F 0x9B 0x98..0x9F  #E0.0   [8] (🛘..🛟)    <reserved-1F6D8>..<...
      | 0xF0 0x9F 0x9B 0xA0..0xA5  #E0.7   [6] (🛠️..🛥️)    hammer and wrench...
      | 0xF0 0x9F 0x9B 0xA6..0xA8  #E0.0   [3] (🛦..🛨)    UP-POINTING MILITAR...
      | 0xF0 0x9F 0x9B 0xA9     #E0.7   [1] (🛩️)       small airplane
      | 0xF0 0x9F 0x9B 0xAA     #E0.0   [1] (🛪)       NORTHEAST-POINTING AIR...
      | 0xF0 0x9F 0x9B 0xAB..0xAC  #E1.0   [2] (🛫..🛬)    airplane departure....
      | 0xF0 0x9F 0x9B 0xAD..0xAF  #E0.0   [3] (🛭..🛯)    <reserved-1F6ED>..<...
      | 0xF0 0x9F 0x9B 0xB0     #E0.7   [1] (🛰️)       satellite
      | 0xF0 0x9F 0x9B 0xB1..0xB2  #E0.0   [2] (🛱..🛲)    ONCOMING FIRE ENGIN...
      | 0xF0 0x9F 0x9B 0xB3     #E0.7   [1] (🛳️)       passenger ship
      | 0xF0 0x9F 0x9B 0xB4..0xB6  #E3.0   [3] (🛴..🛶)    kick scooter..canoe
      | 0xF0 0x9F 0x9B 0xB7..0xB8  #E5.0   [2] (🛷..🛸)    sled..flying saucer
      | 0xF0 0x9F 0x9B 0xB9     #E11.0  [1] (🛹)       skateboard
      | 0xF0 0x9F 0x9B 0xBA     #E12.0  [1] (🛺)       auto rickshaw
      | 0xF0 0x9F 0x9B 0xBB..0xBC  #E13.0  [2] (🛻..🛼)    pickup truck..rolle...
      | 0xF0 0x9F 0x9B 0xBD..0xBF  #E0.0   [3] (🛽..🛿)    <reserved-1F6FD>..<...
      | 0xF0 0x9F 0x9D 0xB4..0xBF  #E0.0  [12] (🝴..🝿)    <reserved-1F774>..<...
      | 0xF0 0x9F 0x9F 0x95..0x9F  #E0.0  [11] (🟕..🟟)    CIRCLED TRIANGLE..<...
      | 0xF0 0x9F 0x9F 0xA0..0xAB  #E12.0 [12] (🟠..🟫)    orange circle..brow...
      | 0xF0 0x9F 0x9F 0xAC..0xBF  #E0.0  [20] (🟬..🟿)    <reserved-1F7EC>..<...
      | 0xF0 0x9F 0xA0 0x8C..0x8F  #E0.0   [4] (🠌..🠏)    <reserved-1F80C>..<...
      | 0xF0 0x9F 0xA1 0x88..0x8F  #E0.0   [8] (🡈..🡏)    <reserved-1F848>..<...
      | 0xF0 0x9F 0xA1 0x9A..0x9F  #E0.0   [6] (🡚..🡟)    <reserved-1F85A>..<...
      | 0xF0 0x9F 0xA2 0x88..0x8F  #E0.0   [8] (🢈..🢏)    <reserved-1F888>..<...
      | 0xF0 0x9F 0xA2 0xAE..0xFF  #E0.0  [82] (🢮..🣿)    <reserved-1F8AE>..<...
      | 0xF0 0x9F 0xA3 0x00..0xBF  #
      | 0xF0 0x9F 0xA4 0x8C     #E13.0  [1] (🤌)       pinched fingers
      | 0xF0 0x9F 0xA4 0x8D..0x8F  #E12.0  [3] (🤍..🤏)    white heart..pinchi...
      | 0xF0 0x9F 0xA4 0x90..0x98  #E1.0   [9] (🤐..🤘)    zipper-mouth face.....
      | 0xF0 0x9F 0xA4 0x99..0x9E  #E3.0   [6] (🤙..🤞)    call me hand..cross...
      | 0xF0 0x9F 0xA4 0x9F     #E5.0   [1] (🤟)       love-you gesture
      | 0xF0 0x9F 0xA4 0xA0..0xA7  #E3.0   [8] (🤠..🤧)    cowboy hat face..sn...
      | 0xF0 0x9F 0xA4 0xA8..0xAF  #E5.0   [8] (🤨..🤯)    face with raised ey...
      | 0xF0 0x9F 0xA4 0xB0     #E3.0   [1] (🤰)       pregnant woman
      | 0xF0 0x9F 0xA4 0xB1..0xB2  #E5.0   [2] (🤱..🤲)    breast-feeding..pal...
      | 0xF0 0x9F 0xA4 0xB3..0xBA  #E3.0   [8] (🤳..🤺)    selfie..person fencing
      | 0xF0 0x9F 0xA4 0xBC..0xBE  #E3.0   [3] (🤼..🤾)    people wrestling..p...
      | 0xF0 0x9F 0xA4 0xBF     #E12.0  [1] (🤿)       diving mask
      | 0xF0 0x9F 0xA5 0x80..0x85  #E3.0   [6] (🥀..🥅)    wilted flower..goal...
      | 0xF0 0x9F 0xA5 0x87..0x8B  #E3.0   [5] (🥇..🥋)    1st place medal..ma...
      | 0xF0 0x9F 0xA5 0x8C     #E5.0   [1] (🥌)       curling stone
      | 0xF0 0x9F 0xA5 0x8D..0x8F  #E11.0  [3] (🥍..🥏)    lacrosse..flying disc
      | 0xF0 0x9F 0xA5 0x90..0x9E  #E3.0  [15] (🥐..🥞)    croissant..pancakes
      | 0xF0 0x9F 0xA5 0x9F..0xAB  #E5.0  [13] (🥟..🥫)    dumpling..canned food
      | 0xF0 0x9F 0xA5 0xAC..0xB0  #E11.0  [5] (🥬..🥰)    leafy green..smilin...
      | 0xF0 0x9F 0xA5 0xB1     #E12.0  [1] (🥱)       yawning face
      | 0xF0 0x9F 0xA5 0xB2     #E13.0  [1] (🥲)       smiling face with tear
      | 0xF0 0x9F 0xA5 0xB3..0xB6  #E11.0  [4] (🥳..🥶)    partying face..cold...
      | 0xF0 0x9F 0xA5 0xB7..0xB8  #E13.0  [2] (🥷..🥸)    ninja..disguised face
      | 0xF0 0x9F 0xA5 0xB9     #E0.0   [1] (🥹)       <reserved-1F979>
      | 0xF0 0x9F 0xA5 0xBA     #E11.0  [1] (🥺)       pleading face
      | 0xF0 0x9F 0xA5 0xBB     #E12.0  [1] (🥻)       sari
      | 0xF0 0x9F 0xA5 0xBC..0xBF  #E11.0  [4] (🥼..🥿)    lab coat..flat shoe
      | 0xF0 0x9F 0xA6 0x80..0x84  #E1.0   [5] (🦀..🦄)    crab..unicorn
      | 0xF0 0x9F 0xA6 0x85..0x91  #E3.0  [13] (🦅..🦑)    eagle..squid
      | 0xF0 0x9F 0xA6 0x92..0x97  #E5.0   [6] (🦒..🦗)    giraffe..cricket
      | 0xF0 0x9F 0xA6 0x98..0xA2  #E11.0 [11] (🦘..🦢)    kangaroo..swan
      | 0xF0 0x9F 0xA6 0xA3..0xA4  #E13.0  [2] (🦣..🦤)    mammoth..dodo
      | 0xF0 0x9F 0xA6 0xA5..0xAA  #E12.0  [6] (🦥..🦪)    sloth..oyster
      | 0xF0 0x9F 0xA6 0xAB..0xAD  #E13.0  [3] (🦫..🦭)    beaver..seal
      | 0xF0 0x9F 0xA6 0xAE..0xAF  #E12.0  [2] (🦮..🦯)    guide dog..white cane
      | 0xF0 0x9F 0xA6 0xB0..0xB9  #E11.0 [10] (🦰..🦹)    red hair..supervillain
      | 0xF0 0x9F 0xA6 0xBA..0xBF  #E12.0  [6] (🦺..🦿)    safety vest..mechan...
      | 0xF0 0x9F 0xA7 0x80     #E1.0   [1] (🧀)       cheese wedge
      | 0xF0 0x9F 0xA7 0x81..0x82  #E11.0  [2] (🧁..🧂)    cupcake..salt
      | 0xF0 0x9F 0xA7 0x83..0x8A  #E12.0  [8] (🧃..🧊)    beverage box..ice
      | 0xF0 0x9F 0xA7 0x8B     #E13.0  [1] (🧋)       bubble tea
      | 0xF0 0x9F 0xA7 0x8C     #E0.0   [1] (🧌)       <reserved-1F9CC>
      | 0xF0 0x9F 0xA7 0x8D..0x8F  #E12.0  [3] (🧍..🧏)    person standing..de...
      | 0xF0 0x9F 0xA7 0x90..0xA6  #E5.0  [23] (🧐..🧦)    face with monocle.....
      | 0xF0 0x9F 0xA7 0xA7..0xBF  #E11.0 [25] (🧧..🧿)    red envelope..nazar...
      | 0xF0 0x9F 0xA8 0x80..0xFF  #E0.0 [112] (🨀..🩯)    NEUTRAL CHESS KING....
      | 0xF0 0x9F 0xA9 0x00..0xAF  #
      | 0xF0 0x9F 0xA9 0xB0..0xB3  #E12.0  [4] (🩰..🩳)    ballet shoes..shorts
      | 0xF0 0x9F 0xA9 0xB4     #E13.0  [1] (🩴)       thong sandal
      | 0xF0 0x9F 0xA9 0xB5..0xB7  #E0.0   [3] (🩵..🩷)    <reserved-1FA75>..<...
      | 0xF0 0x9F 0xA9 0xB8..0xBA  #E12.0  [3] (🩸..🩺)    drop of blood..stet...
      | 0xF0 0x9F 0xA9 0xBB..0xBF  #E0.0   [5] (🩻..🩿)    <reserved-1FA7B>..<...
      | 0xF0 0x9F 0xAA 0x80..0x82  #E12.0  [3] (🪀..🪂)    yo-yo..parachute
      | 0xF0 0x9F 0xAA 0x83..0x86  #E13.0  [4] (🪃..🪆)    boomerang..nesting ...
      | 0xF0 0x9F 0xAA 0x87..0x8F  #E0.0   [9] (🪇..🪏)    <reserved-1FA87>..<...
      | 0xF0 0x9F 0xAA 0x90..0x95  #E12.0  [6] (🪐..🪕)    ringed planet..banjo
      | 0xF0 0x9F 0xAA 0x96..0xA8  #E13.0 [19] (🪖..🪨)    military helmet..rock
      | 0xF0 0x9F 0xAA 0xA9..0xAF  #E0.0   [7] (🪩..🪯)    <reserved-1FAA9>..<...
      | 0xF0 0x9F 0xAA 0xB0..0xB6  #E13.0  [7] (🪰..🪶)    fly..feather
      | 0xF0 0x9F 0xAA 0xB7..0xBF  #E0.0   [9] (🪷..🪿)    <reserved-1FAB7>..<...
      | 0xF0 0x9F 0xAB 0x80..0x82  #E13.0  [3] (🫀..🫂)    anatomical heart..p...
      | 0xF0 0x9F 0xAB 0x83..0x8F  #E0.0  [13] (🫃..🫏)    <reserved-1FAC3>..<...
      | 0xF0 0x9F 0xAB 0x90..0x96  #E13.0  [7] (🫐..🫖)    blueberries..teapot
      | 0xF0 0x9F 0xAB 0x97..0xBF  #E0.0  [41] (🫗..🫿)    <reserved-1FAD7>..<...
      | 0xF0 0x9F 0xB0 0x80..0xFF        #E0.0[1022] (🰀..🿽)    <reserved-1FC...
      | 0xF0 0x9F 0xB1..0xBE 0x00..0xFF  #
      | 0xF0 0x9F 0xBF 0x00..0xBD        #
      ;

}%%
