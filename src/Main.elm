module Main exposing (main)

import Base64
import Browser
import Bytes.Encode as Encode
import Flate
import Html exposing (Html)
import Html.Attributes as Attr
import Html.Events as Events
import Url



-- MAIN


main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }



-- MODEL


type alias Model =
    { language : String
    , code : String
    }


init : Model
init =
    { language = "JavaScript"
    , code = "const foo = \"bar\""
    }



-- UPDATE


type Msg
    = TypedLanguage String
    | TypedCode String


update : Msg -> Model -> Model
update msg model =
    case msg of
        TypedLanguage l ->
            { model | language = l }

        TypedCode c ->
            { model | code = c }



-- VIEW


view : Model -> Html Msg
view model =
    Html.div [ Attr.class "monospace h-screen w-screen relative" ]
        [ Html.input
            [ Attr.class "fixed p-1 border-solid border-2 bg-black text-white rounded-xl top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-center"
            , Attr.value model.language
            , Attr.title "Select language for syntax highlighting"
            , Events.onInput TypedLanguage
            ]
            []
        , Html.div [ Attr.class "fixed bottom-0 right-0 text-center text-sm p-2" ]
            [ Html.a [ Attr.href "https://enso.no", Attr.target "_blank" ] [ Html.text "Made with ❤️ by Ensō" ] ]
        , Html.textarea
            [ halfScreenAttr
            , Attr.class "outline-none border-solid border-2 rounded-xl p-4 font-mono resize-none"
            , Attr.value model.code
            , Events.onInput TypedCode
            , Attr.autofocus True
            ]
            []
        , Html.div [ halfScreenAttr, Attr.class "flex items-center justify-center" ]
            [ Html.img
                [ Attr.src <| codeToSrc model
                ]
                []
            ]
        ]


halfScreenAttr : Html.Attribute msg
halfScreenAttr =
    Attr.class "float-left md:h-full md:w-1/2 h-1/2 w-full overflow-y-scroll p-4"


codeToSrc : { a | code : String, language : String } -> String
codeToSrc { code, language } =
    --"https://codimg.alwaysdata.net/code.svg?input="
    "/code.svg?input="
        ++ encodeCodeBlock code
        ++ "&lang="
        ++ language


encodeCodeBlock : String -> String
encodeCodeBlock =
    Encode.string
        >> Encode.encode
        >> Flate.deflate
        >> Base64.fromBytes
        >> Maybe.withDefault "invalid data"
        >> Url.percentEncode
